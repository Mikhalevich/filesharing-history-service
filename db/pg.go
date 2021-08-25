package db

import (
	"database/sql"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connectionStr string) (*Postgres, error) {
	pgDB, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db: pgDB,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) StoreEvent(e *Event) error {
	return nil
}
