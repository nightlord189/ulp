package model

import (
	"errors"
	"time"
)

const (
	RoleStudent Role = "student"
	RoleTutor   Role = "tutor"
)

type Role string

type UserDB struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"column:username;unique"`
	Email        string    `json:"email" gorm:"column:email"`
	Role         Role      `json:"type" gorm:"column:role"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (UserDB) TableName() string {
	return "users"
}

type RegRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Email    string `json:"email" binding:"required" form:"email"`
	Role     string `json:"role" binding:"required" form:"role"`
	Password string `json:"password" binding:"required" form:"password"`
}

func (r *RegRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	if len(r.Password) < 6 {
		return errors.New("length of password should be at least 6 symbols")
	}
	if r.Email == "" {
		return errors.New("email is empty")
	}
	if r.Role != string(RoleStudent) && r.Role != string(RoleTutor) {
		return errors.New("invalid role")
	}
	return nil
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
