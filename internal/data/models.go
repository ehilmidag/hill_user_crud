package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("UserEntity with that id does not exist")
)

type Models struct {
	Users UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
