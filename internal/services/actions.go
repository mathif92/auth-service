package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	errs "github.com/mathif92/auth-service/internal/errors"
	"github.com/pkg/errors"
)

const (
	_insertAction = `INSERT INTO action (action, entity, updated_at) 
								VALUES (?, ?, current_timestamp())`

	_selectAction = `
		SELECT
			a.id,
			a.action,
			a.entity,
			a.enabled,
			a.created_at,
			a.updated_at
		FROM
			action as a
		WHERE a.id = ?
	`
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

func (a *Actions) GetAction(ctx context.Context, actionID int64) (ActionModel, error) {
	var action ActionModel
	if err := a.db.GetContext(ctx, &action, _selectAction, actionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ActionModel{}, errs.New("action not found", http.StatusNotFound)
		}
		return ActionModel{}, err
	}

	return action, nil
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

	queryArgs["updated_at"] = goqu.L("current_timestamp()")

	dialect := goqu.Dialect("mysql")
	updateQuery, _, err := dialect.Update("action").Set(
		goqu.Record(queryArgs),
	).Where(goqu.Ex{"id": goqu.Op{"eq": action.ID}}).ToSQL()

	if err != nil {
		return errors.Wrap(err, "creating update action db query")
	}

	if _, err := tx.ExecContext(ctx, updateQuery); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating action in db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}
