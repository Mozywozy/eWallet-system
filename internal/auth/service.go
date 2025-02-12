package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(user User) (*User, error)
}

type authService struct {
	userRepo UserRepository 
}

func NewAuthService(repo UserRepository) AuthService { 
	return &authService{userRepo: repo}
}

func (s *authService) RegisterUser(user User) (*User, error) {
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
