package model

type RequestType string

var (
	RequestTypeIn  RequestType = "in"
	RequestTypeOut RequestType = "out"
)

type Request struct {
	ID       int         `db:"id"`
	Type     RequestType `db:"request_type"`
	TgUserID int         `db:"tg_user_id"`
	Tag      Tag         `db:"tag"`

	TagID int `db:"_"`
}

type RequestQueryOptions struct {
	Type  []RequestType
	TagID []int
}
