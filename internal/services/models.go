package services

import (
	"time"
)

type CredentialsModel struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type TokenModel struct {
	ID            int64     `db:"id"`
	Token         string    `db:"token"`
	CredentialsID int64     `db:"credentials_id"`
	TimeToLiveAt  time.Time `db:"ttl_at"`
	CreatedAt     time.Time `db:"created_at"`
}

type RoleModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Enabled   bool      `db:"enabled"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ActionModel struct {
	ID        int64     `db:"id"`
	Action    string    `db:"action"`
	Entity    string    `db:"entity"`
	Enabled   bool      `db:"enabled"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RolesActionsModel struct {
	ID       int64 `db:"id"`
	RoleID   int64 `db:"role_id"`
	ActionID int64 `db:"action_id"`
}

type RolesCredentialsModel struct {
	ID            int64 `db:"id"`
	RoleID        int64 `db:"role_id"`
	CredentialsID int64 `db:"credentials_id"`
}
