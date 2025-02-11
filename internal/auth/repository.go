package auth

import (
	"ewallet-engine/internal/database"

	"gorm.io/gorm"
)

// UserRepository interface untuk manipulasi data user
type UserRepository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
}

// userRepository struct
type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository membuat instance repository baru
func NewUserRepository(dbService database.Service) UserRepository { // 🔥 Kembalikan sebagai UserRepository
	return &userRepository{DB: dbService.GetDB()}
}

// CreateUser menyimpan user ke database
func (r *userRepository) CreateUser(user *User) error {
	return r.DB.Create(user).Error
}

// FindByEmail mencari user berdasarkan email
func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
