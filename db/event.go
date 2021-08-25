package db

const (
	Add = iota + 1
	Remove
)

type Event struct {
	UserID   int64
	UserName string
	FileName string
	Time     int64
	Size     int64
	Action   int
}