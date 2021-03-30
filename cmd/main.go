package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mehrdod/todo"
	"github.com/mehrdod/todo/pkg/handler"
	"github.com/mehrdod/todo/pkg/repository"
	"github.com/mehrdod/todo/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error init configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error, failed to load env: %s", err.Error())
	}

	srv := new(todo.Server)

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error, failed init db %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	if err := srv.Startup(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error during the server startup: %s", err.Error())
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("conf")
	return viper.ReadInConfig()
}
