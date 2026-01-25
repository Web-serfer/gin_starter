package repository

import (
	"gin-starter/internal/models"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}
