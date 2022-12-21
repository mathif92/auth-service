package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mathif92/auth-service/internal/db"
	"github.com/pkg/errors"
)

type Actions struct {
	db *sqlx.DB
}

func NewActions(db *sqlx.DB) *Actions {
	return &Actions{db: db}
}

func (a *Actions) SaveAction(ctx context.Context, action ActionModel) (int64, error) {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, "creating db transacion")
	}
	const (
		_insertAction = `INSERT INTO action (action, entity) 
								VALUES (?, ?)`
	)

	result, err := tx.ExecContext(ctx, _insertAction, action.Action, action.Entity)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "inserting action into db")
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "commiting db transaction")
	}

	return result.LastInsertId()
}

func (a *Actions) UpdateAction(ctx context.Context, action ActionModel, enabledFlag *bool) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}

	queryArgs := make(map[string]interface{})
	if action.Action != "" {
		queryArgs["action"] = action.Action
	}

	if action.Entity != "" {
		queryArgs["entity"] = action.Entity
	}

	if enabledFlag != nil {
		queryArgs["enabled"] = *enabledFlag
	}

	queryArgs["updated_at"] = "current_timestamp()"

	if err := db.Update(ctx, tx, "action", queryArgs, map[string]interface{}{"id": action.ID}); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating action in db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}
