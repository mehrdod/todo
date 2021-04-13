package service

import (
	"github.com/mehrdod/todo/domain"
	"github.com/mehrdod/todo/pkg/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem) *TodoItemService {
	return &TodoItemService{repo: repo}
}

func (ts *TodoItemService) Create(userId int, list domain.TodoItem) (int, error) {
	return ts.repo.Create(userId, list)
}
