package services

import "github.com/jmoiron/sqlx"

type Actions struct {
	db *sqlx.DB
}

func NewActions(db *sqlx.DB) *Actions {
	return &Actions{db: db}
}
