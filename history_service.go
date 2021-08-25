package main

import (
	"context"

	"github.com/Mikhalevich/filesharing-history-service/db"
	"github.com/Mikhalevich/filesharing/httpcode"
	"github.com/Mikhalevich/filesharing/proto/event"
	"github.com/Mikhalevich/filesharing/proto/history"
)

type storager interface {
	StoreEvent(*db.Event) error
	EventsByUserID(userID int64) ([]*db.Event, error)
}

type HistoryService struct {
	storage storager
}

func NewHistoryService(s storager) *HistoryService {
	return &HistoryService{
		storage: s,
	}
}

func marshalFileEvent(e *db.Event) *event.FileEvent {
	return &event.FileEvent{
		UserID:   e.UserID,
		UserName: e.UserName,
		FileName: e.FileName,
		Time:     e.Time,
		Size:     e.Size,
		Action:   event.Action(e.Action),
	}
}

func unmarshalFileEvent(e *event.FileEvent) *db.Event {
	return &db.Event{
		UserID:   e.GetUserID(),
		UserName: e.GetUserName(),
		FileName: e.GetFileName(),
		Time:     e.GetTime(),
		Size:     e.GetSize(),
		Action:   int(e.GetAction()),
	}
}

func (hs *HistoryService) List(ctx context.Context, req *history.ListRequest, rsp *history.ListResponse) error {
	events, err := hs.storage.EventsByUserID(req.GetUserID())
	if err != nil {
		return httpcode.NewWrapInternalServerError(err, "unable to get events by user id")
	}

	for _, e := range events {
		rsp.Files = append(rsp.Files, marshalFileEvent(e))
	}
	return nil
}

func (hs *HistoryService) StoreEvent(ctx context.Context, req *event.FileEvent) error {
	if err := hs.storage.StoreEvent(unmarshalFileEvent(req)); err != nil {
		return httpcode.NewWrapInternalServerError(err, "unable to store file event")
	}
	return nil
}
