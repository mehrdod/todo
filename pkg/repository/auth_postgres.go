package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mehrdod/todo/domain"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

const createUserSql = `
	INSERT INTO %s (name, username, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id`

func (ap *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf(createUserSql, userTable)
	row := ap.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const getUserSql = `
	SELECT id FROM %s
	WHERE username = $1 AND password_hash = $2`

func (ap *AuthPostgres) GetUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(getUserSql, userTable)
	err := ap.db.Get(&user, query, username, password)

	return user, err
}
