package model

import (
	"errors"
	"time"
)

type UserDB struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"column:username"`
	Email        string    `json:"email" gorm:"column:email"`
	Type         string    `json:"type" gorm:"column:type"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (UserDB) TableName() string {
	return "users"
}

type RegRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

func (r *AuthRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
