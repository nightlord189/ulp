package service

import (
	"errors"
	"fmt"
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"github.com/nightlord189/ulp/internal/model"
	"gorm.io/gorm"
	"time"
)

// Service - структура со ссылками на зависимости
type Service struct {
	Config *config.Config
	DB     *db.Manager
}

// NewService - конструктор Service
func NewService(cfg *config.Config, db *db.Manager) *Service {
	service := Service{
		Config: cfg,
		DB:     db,
	}
	return &service
}

func (s *Service) Auth(req model.AuthRequest) (string, error) {
	var user model.UserDB
	err := s.DB.GetEntityByField("username", req.Username, &user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("invalid username or password")
		}
		return "", fmt.Errorf("err get user from DB: %w", err)
	}
	if GetHash(req.Password) != user.PasswordHash {
		return "", errors.New("invalid username or password")
	}
	payload := CreatePayload(user.ID, user.Username, user.Type,
		time.Now().Add(time.Second*time.Duration(s.Config.Auth.ExpTime)).Unix(), s.Config.Auth.Issuer)
	token, err := CreateToken(payload, s.Config.Auth.Secret)
	if err != nil {
		return "", fmt.Errorf("err on create token: %w", err)
	}
	return token, nil
}
