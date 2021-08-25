package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Mikhalevich/filesharing-history-service/db"
	"github.com/Mikhalevich/filesharing/proto/history"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/sirupsen/logrus"
)

type params struct {
	ServiceName        string
	DBConnectionString string
}

func loadParams() (*params, error) {
	var p params
	p.ServiceName = os.Getenv("FS_SERVICE_NAME")
	if p.ServiceName == "" {
		p.ServiceName = "history.service"
	}

	p.DBConnectionString = os.Getenv("FS_DB_CONNECTION_STRING")
	if p.DBConnectionString == "" {
		return nil, errors.New("databse connection string is missing, please specify FS_DB_CONNECTION_STRING environment variable")
	}

	return &p, nil
}

func makeLoggerWrapper(logger *logrus.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			logger.Infof("processing %s", req.Method())
			start := time.Now()
			defer logger.Infof("end processing %s, time = %v", req.Method(), time.Now().Sub(start))
			err := fn(ctx, req, rsp)
			if err != nil {
				logger.Errorln(err)
			}
			return err
		}
	}
}

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	p, err := loadParams()
	if err != nil {
		logger.Errorln(fmt.Errorf("unable to load params: %W", err))
		return
	}

	logger.Infof("running auth service with params: %v\n", p)

	srv := micro.NewService(
		micro.Name(p.ServiceName),
		micro.WrapHandler(makeLoggerWrapper(logger)),
	)

	srv.Init()

	var storage *db.Postgres
	for i := 0; i < 3; i++ {
		storage, err = db.NewPostgres(p.DBConnectionString)
		if err == nil {
			break
		}

		time.Sleep(time.Second * 1)
		logger.Infof("try to connect to database: %d  error: %v\n", i, err)
	}

	if err != nil {
		logger.Errorln(fmt.Errorf("unable to connect to database: %w", err))
		return
	}
	defer storage.Close()

	hs := NewHistoryService(storage)

	micro.RegisterSubscriber("filesharing.file.event", srv.Server(), hs.StoreEvent, server.SubscriberQueue("filesharing.history.service.queue"))
	history.RegisterHistoryServiceHandler(srv.Server(), hs)

	err = srv.Run()
	if err != nil {
		logger.Errorln(err)
		return
	}
}
