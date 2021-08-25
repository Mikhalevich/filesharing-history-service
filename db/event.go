package db

type Event struct {
	UserID   int64  `db:"user_id"`
	UserName string `db:"user_name"`
	FileName string `db:"file_name"`
	Time     int64  `db:"time"`
	Size     int64  `db:"size"`
	Action   int    `db:"action"`
}
