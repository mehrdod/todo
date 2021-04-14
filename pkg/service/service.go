package service

import (
	"github.com/mehrdod/todo/domain"
	"github.com/mehrdod/todo/pkg/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId int, listId int) (domain.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, request domain.UpdateListRequest) error
}

type TodoItem interface {
	Create(userId int, listId int, list domain.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]domain.TodoItem, error)
	GetById(userId int, itemId int) (domain.TodoItem, error)
	Delete(userId int, itemId int) error
	Update(userId int, listId int, request domain.UpdateItemRequest) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
