package web

import "github.com/google/uuid"

type UserID string

func (u UserID) ToUUID() (ID uuid.UUID, err error) {
	ID, err = uuid.Parse(string(u))
	return
}

type Params struct {
	UserID UserID `validate:"required"`
	Query  interface{}
}

func NewParams(userID string, query interface{}) *Params {
	return &Params{UserID: UserID(userID), Query: query}
}
