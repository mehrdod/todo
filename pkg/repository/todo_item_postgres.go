package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mehrdod/todo/domain"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

const (
	createItemSql = `
	INSERT INTO %s (title, description)
	VALUES ($1, $2)
	RETURNING id
`
	createListsItemSql = `
	INSERT INTO %s (item_id, list_id)
	VALUES ($1, $2)
`
)

func (tr *TodoItemPostgres) Create(listId int, item domain.TodoItem) (int, error) {
	tx, err := tr.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createItemQuery := fmt.Sprintf(createItemSql, todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf(createListsItemSql, userListsTable)
	_, err = tx.Exec(createListsItemsQuery, listId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
