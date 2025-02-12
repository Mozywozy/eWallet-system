package auth

import (
	"context"
	"errors"
	"ewallet-engine/internal/database"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	SaveUserSession(session *UserSession) error
	SaveTokenToCache(userID uint, token, refreshToken string, expiration time.Duration) error
	DeleteTokenFromCache(userID uint, refreshToken string) error
	FindSessionByRefreshToken(token string) (*UserSession, error)
	DeleteUserSession(userID uint) error
	FindUserIDByRefreshToken(refreshToken string) (uint, error)
	GetRedis() *redis.Client
}

type userRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (r *userRepository) GetRedis() *redis.Client {
	return r.Redis
}

func NewUserRepository(dbService database.Service) UserRepository {
	return &userRepository{DB: dbService.GetDB(),
		Redis: dbService.GetRedis(),
	}
}

func (r *userRepository) SaveTokenToCache(userID uint, token, refreshToken string, expiration time.Duration) error {
	ctx := context.Background()
	pipe := r.Redis.Pipeline()

	pipe.Set(ctx, getTokenKey(userID), token, expiration)

	pipe.Set(ctx, getRefreshTokenKey(refreshToken), userID, expiration*24)

	_, err := pipe.Exec(ctx)
	return err
}


func (r *userRepository) DeleteTokenFromCache(userID uint, refreshToken string) error {
	ctx := context.Background()
	pipe := r.Redis.TxPipeline()

	pipe.Del(ctx, getTokenKey(userID))

	pipe.Del(ctx, getRefreshTokenKey(refreshToken))

	_, err := pipe.Exec(ctx)
	return err
}


func getTokenKey(userID uint) string {
	return "auth:token:" + fmt.Sprint(userID)
}

func getRefreshTokenKey(refreshToken string) string {
	return "auth:refresh:" + refreshToken
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

func (r *userRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SaveUserSession(session *UserSession) error {
	return r.DB.Create(session).Error
}

func (r *userRepository) DeleteUserSession(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&UserSession{}).Error
}

func (r *userRepository) FindSessionByRefreshToken(token string) (*UserSession, error) {
	var session UserSession
	err := r.DB.Where("refresh_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *userRepository) FindUserIDByRefreshToken(refreshToken string) (uint, error) {
	ctx := context.Background()

	exists, err := r.Redis.Exists(ctx, getRefreshTokenKey(refreshToken)).Result()
	if err != nil || exists == 0 {
		return 0, errors.New("refresh token tidak ditemukan di cache")
	}

	userID, err := r.Redis.Get(ctx, getRefreshTokenKey(refreshToken)).Uint64()
	if err != nil {
		return 0, errors.New("gagal mengambil userID dari Redis")
	}

	return uint(userID), nil
}
