package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mathif92/auth-service/internal/db"
	"github.com/pkg/errors"
)

const (
	_insertRole = `INSERT INTO role (name, updated_at) 
								VALUES (?, current_timestamp())`

	_selectRoleWithActions = `
		SELECT
			r.id 		 as roleID,
			r.name 		 as roleName,
			r.enabled 	 as roleEnabled,
			r.created_at as roleCreatedAt,
			r.updated_at as roleUpdatedAt,
			a.id 		 as actionID,
			a.action 	 as actionAction,
			a.entity 	 as actionEntity,
			a.enabled 	 as actionEnabled,
			a.created_at as actionCreatedAt,
			a.updated_at as actionUpdatedAt
		FROM
			role r LEFT JOIN
			roles_actions ra ON ra.role_id = r.id LEFT JOIN
			action a ON a.id = ra.action_id
		WHERE
			r.id = ?
	`

	_insertRoleAction = `INSERT INTO roles_actions (role_id, action_id) 
								VALUES (?, ?)`

	_deleteRoleAction = `DELETE FROM roles_actions WHERE role_id = ? AND action_id = ?`

	_insertRoleCredentials = `INSERT INTO roles_credentials (role_id, credentials_id) 
								VALUES (?, ?)`

	_deleteRoleCredentials = `DELETE FROM roles_credentials WHERE role_id = ? AND credentials_id = ?`
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

// GetRole returns an existing role from the DB based on the provided roleID, including the assigned actions. Returns an error in case there is one.
func (r *Roles) GetRole(ctx context.Context, roleID int64) (RoleWithActions, error) {
	var rows []roleWithActions
	err := r.db.SelectContext(ctx, &rows, _selectRoleWithActions, roleID)
	if err != nil {
		return RoleWithActions{}, err
	}
	result := RoleWithActions{}
	for _, row := range rows {
		if result.Role.ID == 0 {
			result.Role = RoleModel{
				ID:        row.RoleID,
				Name:      row.RoleName,
				Enabled:   row.RoleEnabled,
				CreatedAt: row.RoleCreatedAt,
				UpdatedAt: row.RoleUpdatedAt,
			}
		}

		if row.ActionID.Int64 > 0 {
			result.Actions = append(result.Actions, ActionModel{
				ID:        row.ActionID.Int64,
				Action:    row.ActionAction.String,
				Entity:    row.ActionEntity.String,
				Enabled:   row.ActionEnabled.Bool,
				CreatedAt: row.ActionCreatedAt.Time,
				UpdatedAt: row.ActionUpdatedAt.Time,
			})
		}
	}

	return result, nil
}

// UpdateRole updates existing role into the DB. Returns an error in case there is one.
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

	_, err = tx.ExecContext(ctx, _insertRoleCredentials, roleCreds.RoleID, roleCreds.CredentialsID)
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

	_, err = tx.ExecContext(ctx, _deleteRoleCredentials, roleCreds.RoleID, roleCreds.CredentialsID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing role from credentials in the db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}
