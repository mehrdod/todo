package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mehrdod/todo/domain"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

const (
	createListSql = `
	INSERT INTO %s (title, description)
	VALUES ($1, $2)
	RETURNING id
`
	createUsersListSql = `
	INSERT INTO %s (user_id, list_id)
	VALUES ($1, $2)
`
)

func (tr *TodoListPostgres) Create(userId int, list domain.TodoList) (int, error) {
	tx, err := tr.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf(createListSql, todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf(createUsersListSql, userListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

const getAllSql = `
	SELECT td.title, td.description FROM %s AS td
	INNER JOIN %s AS ul
	ON ul.list_id = td.id
	WHERE ul.user_id = $1
`

func (tr *TodoListPostgres) GetAll(userId int) ([]domain.TodoList, error) {

	var todoLists []domain.TodoList

	getAllQuery := fmt.Sprintf(getAllSql, todoListsTable, userListsTable)
	err := tr.db.Select(&todoLists, getAllQuery, userId)
	return todoLists, err
}

const getByIdSql = `
	SELECT td.title, td.description FROM %s AS td
	INNER JOIN %s AS ul
	ON ul.list_id = td.id
	WHERE ul.user_id = $1 AND ul.list_id = $2
`

func (tr *TodoListPostgres) GetById(userId int, listId int) (domain.TodoList, error) {

	var list domain.TodoList

	getByIdQuery := fmt.Sprintf(getByIdSql, todoListsTable, userListsTable)
	err := tr.db.Get(&list, getByIdQuery, userId, listId)
	return list, err
}
