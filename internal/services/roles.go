package services

import "github.com/jmoiron/sqlx"

type Roles struct {
	db *sqlx.DB
}

func NewRoles(db *sqlx.DB) *Roles {
	return &Roles{db: db}
}
