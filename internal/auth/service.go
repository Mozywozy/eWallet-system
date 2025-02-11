package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService interface untuk layanan otentikasi
type AuthService interface {
	RegisterUser(user User) (*User, error)
}

// authService struct
type authService struct {
	userRepo UserRepository // 🔥 Gunakan Interface, Bukan Struct
}

// NewAuthService membuat instance service baru
func NewAuthService(repo UserRepository) AuthService { // 🔥 Parameter harus interface
	return &authService{userRepo: repo}
}

// RegisterUser memproses pendaftaran user
func (s *authService) RegisterUser(user User) (*User, error) {
	// Cek apakah email sudah digunakan
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Simpan user ke database
	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
