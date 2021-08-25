package db

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(connectionStr string) (*Postgres, error) {
	pgDB, err := sqlx.Connect("postgres", connectionStr)
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
	const query = `
		INSERT INTO FileEvent(user_id, user_name, file_name, time, size, action) 
		VALUES(:user_id, :user_name, :file_name, :time, :size, :action)
	`
	if _, err := p.db.NamedExec(query, e); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) EventsByUserID(userID int64) ([]*Event, error) {
	const query = `
	SELECT * 
	FROM FileEvent
	WHERE user_id = $1`

	var events []*Event
	if err := p.db.Select(&events, query, userID); err != nil {
		return nil, err
	}

	return events, nil
}
