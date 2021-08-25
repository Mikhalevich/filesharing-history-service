package main

import (
	"context"

	"github.com/Mikhalevich/filesharing-history-service/db"
	"github.com/Mikhalevich/filesharing/proto/history"
)

type storager interface {
	StoreEvent(*db.Event) error
}

type HistoryService struct {
	storage storager
}

func NewHistoryService(s storager) *HistoryService {
	return &HistoryService{
		storage: s,
	}
}

func (hs *HistoryService) List(ctx context.Context, req *history.ListRequest, rsp *history.ListResponse) error {
	return nil
}
