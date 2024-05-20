package psql

import (
	"context"
	"crud-go/internal/entity"
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(ctx context.Context, user entity.User) error {
	_, err := u.db.Exec("INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (u *Users) GetByCredentials(ctx context.Context, email, password string) (entity.User, error) {
	var user entity.User
	err := u.db.QueryRow("SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}
