package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(user User) (*User, error)
	LoginUser(request LoginRequest) (*User, string, string, error)
	LogoutUser(userID uint) error
	RefreshAccessToken(refreshToken string) (string, string, error)
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

func (s *authService) LoginUser(request LoginRequest) (*User, string, string, error) {
	user, err := s.userRepo.FindByUsername(request.Username)
	if err != nil {
		return nil, "", "", errors.New("username atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, "", "", errors.New("username atau password salah")
	}

	token, refreshToken, err := generateJWT(user)
	if err != nil {
		return nil, "", "", errors.New("gagal membuat token")
	}

	session := UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        time.Now().Add(time.Hour * 1),
		RefreshTokenExpired: time.Now().Add(time.Hour * 24),
	}

	err = s.userRepo.SaveUserSession(&session)
	if err != nil {
		return nil, "", "", errors.New("gagal menyimpan sesi login")
	}

	err = s.userRepo.SaveTokenToCache(user.ID, token, refreshToken, time.Hour)
	if err != nil {
		return nil, "", "", errors.New("gagal menyimpan token ke cache")
	}

	return user, token, refreshToken, nil
}


func generateJWT(user *User) (string, string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func (s *authService) LogoutUser(userID uint) error {
	ctx := context.Background()
	refreshToken, err := s.userRepo.GetRedis().Get(ctx, getTokenKey(userID)).Result()
	if err != nil {
		return errors.New("gagal menemukan refresh token di cache")
	}

	_ = s.userRepo.DeleteUserSession(userID)

	return s.userRepo.DeleteTokenFromCache(userID, refreshToken)
}


func (s *authService) RefreshAccessToken(refreshToken string) (string, string, error) {
	userID, err := s.userRepo.FindUserIDByRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("refresh token tidak valid atau sudah kedaluwarsa")
	}

	token, newRefreshToken, err := generateJWT(&User{ID: userID})
	if err != nil {
		return "", "", errors.New("gagal membuat token baru")
	}

	err = s.userRepo.SaveTokenToCache(userID, token, newRefreshToken, time.Hour)
	if err != nil {
		return "", "", errors.New("gagal menyimpan token ke cache")
	}

	return token, newRefreshToken, nil
}
