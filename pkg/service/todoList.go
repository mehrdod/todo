package service

import (
	"github.com/mehrdod/todo/domain"
	"github.com/mehrdod/todo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (ts *TodoListService) Create(userId int, list domain.TodoList) (int, error) {
	return ts.repo.Create(userId, list)
}

func (ts *TodoListService) GetAll(userId int) ([]domain.TodoList, error) {
	return ts.repo.GetAll(userId)
}

func (ts *TodoListService) GetById(userId int, listId int) (domain.TodoList, error) {
	return ts.repo.GetById(userId, listId)
}
