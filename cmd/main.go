package main

import (
	"github.com/mehrdod/todo"
	"github.com/mehrdod/todo/pkg/handler"
	"github.com/mehrdod/todo/pkg/repository"
	"github.com/mehrdod/todo/pkg/service"
	"log"
)

func main() {
	srv := new(todo.Server)

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	if err := srv.Startup("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error during the server startup: %s", err.Error())
	}
}
