package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mehrdod/todo/domain"
	"github.com/sirupsen/logrus"
	"strings"
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

	createListsItemsQuery := fmt.Sprintf(createListsItemSql, listsItemsTable)
	_, err = tx.Exec(createListsItemsQuery, id, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

const getAllItemsSql = `
	SELECT TI.ID,
		TITLE,
		DESCRIPTION
	FROM %s AS TI
	INNER JOIN %s AS LI ON TI.ID = LI.ITEM_ID
	WHERE LIST_ID = $1
`

func (tr *TodoItemPostgres) GetAll(listId int) ([]domain.TodoItem, error) {
	var items []domain.TodoItem

	getAllItemsQuery := fmt.Sprintf(getAllItemsSql, todoItemsTable, listsItemsTable)
	err := tr.db.Select(&items, getAllItemsQuery, listId)
	return items, err
}

const getItemByIdSql = `
	SELECT TI.ID,
	    TI.TITLE,
		TI.DESCRIPTION,
		TI.DONE
	FROM %s AS TI
	INNER JOIN %s AS LI ON TI.ID = LI.ITEM_ID
	INNER JOIN %s AS UL ON UL.LIST_ID = LI.LIST_ID
	WHERE USER_ID = $1 AND ITEM_ID = $2
`

func (tr *TodoItemPostgres) GetById(userId int, itemId int) (domain.TodoItem, error) {
	var item domain.TodoItem

	getItemByIdQuery := fmt.Sprintf(getItemByIdSql, todoItemsTable, listsItemsTable, userListsTable)

	err := tr.db.Get(&item, getItemByIdQuery, userId, itemId)
	return item, err
}

const deleteByIdSql = `
	DELETE
	FROM %s AS TI USING %s AS LI,
		%s AS UL
	WHERE TI.ID = LI.ITEM_ID
					AND UL.LIST_ID = LI.LIST_ID
					AND UL.USER_ID = $1
`

func (tr *TodoItemPostgres) Delete(userId int, itemId int) error {

	deleteByIdQuery := fmt.Sprintf(deleteByIdSql, todoItemsTable, listsItemsTable, userListsTable)

	_, err := tr.db.Exec(deleteByIdQuery, userId, itemId)
	return err
}

const updateItemSql = `
	UPDATE %s TL
	SET %s
	FROM %s UL
	WHERE TL.ID = UL.LIST_ID
					AND UL.LIST_ID = $%d
					AND UL.USER_ID = $%d
`

func (tr *TodoItemPostgres) Update(userId int, listId int, request domain.UpdateItemRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if request.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *request.Title)
		argId++
	}

	if request.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *request.Description)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	updateQuery := fmt.Sprintf(updateItemSql,
		todoListsTable, setQuery, userListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", updateQuery)
	logrus.Debugf("args: %s", args)

	_, err := tr.db.Exec(updateQuery, args...)
	return err
}
