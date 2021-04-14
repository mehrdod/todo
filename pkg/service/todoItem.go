package service

import (
	"github.com/mehrdod/todo/domain"
	"github.com/mehrdod/todo/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (ts *TodoItemService) Create(userId int, listId int, item domain.TodoItem) (int, error) {
	_, err := ts.listRepo.GetById(userId, listId)
	if err != nil {
		// list does not exists or not belongs to the user
		return 0, err
	}
	return ts.repo.Create(listId, item)
}

func (ts *TodoItemService) GetAll(userId int, listId int) ([]domain.TodoItem, error) {
	_, err := ts.listRepo.GetById(userId, listId)
	if err != nil {
		// list does not exists or not belongs to the user
		return nil, err
	}
	return ts.repo.GetAll(listId)
}

func (ts *TodoItemService) GetById(userId int, itemId int) (domain.TodoItem, error) {
	return ts.repo.GetById(userId, itemId)
}

func (ts *TodoItemService) Delete(userId int, itemId int) error {
	return ts.repo.Delete(userId, itemId)
}

func (ts *TodoItemService) Update(userId int, listId int, request domain.UpdateItemRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return ts.repo.Update(userId, listId, request)
}
