package auth

import (
	"ewallet-engine/internal/database"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(dbService database.Service) UserRepository { 
	return &userRepository{DB: dbService.GetDB()}
}

func (r *userRepository) CreateUser(user *User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
