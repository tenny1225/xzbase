package xzbase

import "time"

type Model struct {
	Id string `json:"id" bson:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" bson:"deletedAt"`
}
type ZError struct {
	Code int64
	Msg  string
}

func (z ZError) Error() string {
	return z.Msg
}
