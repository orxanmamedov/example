package repository

import (
	"example/pkg/storage/mysql"
)

type Repo struct {
	db mysql.MySQL
}

func New(db mysql.MySQL) *Repo {
	return &Repo{db: db}
}
