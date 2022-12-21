package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mathif92/auth-service/internal/db"
	"github.com/pkg/errors"
)

type Roles struct {
	db *sqlx.DB
}

func NewRoles(db *sqlx.DB) *Roles {
	return &Roles{db: db}
}

// SaveRole saves role into the DB. Returns the DB ID or an error instead.
func (r *Roles) SaveRole(ctx context.Context, role RoleModel) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, "creating db transacion")
	}
	const (
		_insertRole = `INSERT INTO roles (name, updated_at) 
								VALUES (?, current_timestamp())`
	)

	result, err := tx.ExecContext(ctx, _insertRole, role.Name)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "inserting role into db")
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "commiting db transaction")
	}

	return result.LastInsertId()
}

func (r *Roles) UpdateRole(ctx context.Context, role RoleModel, enabledFlag *bool) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}

	queryArgs := make(map[string]interface{})
	if role.Name != "" {
		queryArgs["name"] = role.Name
	}

	if enabledFlag != nil {
		queryArgs["enabled"] = *enabledFlag
	}

	queryArgs["updated_at"] = "current_timestamp()"

	if err := db.Update(ctx, tx, "role", queryArgs, map[string]interface{}{"id": role.ID}); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating role in db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}

func (r *Roles) AddActionToRole(ctx context.Context, roleAction RolesActionsModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}
	const (
		_insertRoleAction = `INSERT INTO roles_actions (role_id, action_id) 
								VALUES (?, ?)`
	)

	_, err = tx.ExecContext(ctx, _insertRoleAction, roleAction.RoleID, roleAction.ActionID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "adding action to role in the db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}

func (r *Roles) RemoveActionFromRole(ctx context.Context, roleAction RolesActionsModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}
	const (
		_deleteRoleAction = `DELETE FROM roles_actions WHERE role_id = ? AND action_id = ?`
	)

	_, err = tx.ExecContext(ctx, _deleteRoleAction, roleAction.RoleID, roleAction.ActionID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing action from role in the db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}

func (r *Roles) AddRoleToCredentials(ctx context.Context, roleCreds RolesCredentialsModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}
	const (
		_insertRoleAction = `INSERT INTO roles_credentials (role_id, credentials_id) 
								VALUES (?, ?)`
	)

	_, err = tx.ExecContext(ctx, _insertRoleAction, roleCreds.RoleID, roleCreds.CredentialsID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "adding role to credentials in the db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}

func (r *Roles) UnassignRole(ctx context.Context, roleCreds RolesCredentialsModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}
	const (
		_deleteRoleAction = `DELETE FROM roles_credentials WHERE role_id = ? AND credentials_id = ?`
	)

	_, err = tx.ExecContext(ctx, _deleteRoleAction, roleCreds.RoleID, roleCreds.CredentialsID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing role from credentials in the db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}
