package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/mehrdod/todo/domain"
	"github.com/mehrdod/todo/pkg/repository"
	"os"
)

var salt = os.Getenv("HASH_SALT")

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (as *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return as.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
