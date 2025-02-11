package auth

import (
	"errors"
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	Email       string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(12);unique;not null" json:"phone_number"`
	Address     string    `gorm:"type:text;not null" json:"address"`
	DOB         time.Time `gorm:"type:date;not null" json:"dob"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}


type RegisterRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address" validate:"required"`
	DOB         string `json:"dob" validate:"required"` 
}

func (r *RegisterRequest) ConvertToUser() (*User, error) {
	parsedDOB, err := time.Parse("2006-01-02", r.DOB)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid, gunakan format YYYY-MM-DD")
	}

	return &User{
		Username:    r.Username,
		Password:    r.Password,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		Address:     r.Address,
		DOB:         parsedDOB,
	}, nil
}