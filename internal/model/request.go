package model

type RequestType string

var (
	RequestTypeIn  RequestType = "in"
	RequestTypeOut RequestType = "out"
)

type Request struct {
	ID         int         `db:"id"`
	Type       RequestType `db:"type"`
	TgUsername string      `db:"tg_username"`
	Tag        Tag         `db:"tag"`

	TagID int `db:"_"`
}

type RequestQueryOptions struct {
	TgUsername []string
	Type       []RequestType
	TagID      []int
}
