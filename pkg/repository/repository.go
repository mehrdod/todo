package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mehrdod/todo/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId int, listId int) (domain.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, request domain.UpdateListRequest) error
}

type TodoItem interface {
	Create(listId int, item domain.TodoItem) (int, error)
	GetAll(listId int) ([]domain.TodoItem, error)
	GetById(userId int, itemId int) (domain.TodoItem, error)
	Delete(userId int, itemId int) error
	Update(userId int, listId int, request domain.UpdateItemRequest) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
