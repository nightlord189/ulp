package service

import (
	"errors"
	"fmt"
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"github.com/nightlord189/ulp/internal/model"
	"gorm.io/gorm"
	"strconv"
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
	payload := CreatePayload(user.ID, user.Username, string(user.Role),
		time.Now().Add(time.Second*time.Duration(s.Config.Auth.ExpTime)).Unix(), s.Config.Auth.Issuer)
	token, err := CreateToken(payload, s.Config.Auth.Secret)
	if err != nil {
		return "", fmt.Errorf("err on create token: %w", err)
	}
	return token, nil
}

func (s *Service) Reg(req model.RegRequest) error {
	user := model.UserDB{
		Username:     req.Username,
		Email:        req.Email,
		Role:         model.Role(req.Role),
		PasswordHash: GetHash(req.Password),
	}
	err := s.DB.CreateEntity(&user)
	if err != nil {
		return fmt.Errorf("err create user: %w", err)
	}
	return nil
}

func (s *Service) GetTasks(userID string) (model.TemplateTasks, error) {
	tasks := make([]model.TaskDB, 0)
	err := s.DB.GetEntitiesByField("creator_id", userID, &tasks)
	taskViews := make([]model.TaskView, len(tasks))
	for i := range tasks {
		taskViews[i] = tasks[i].ToView()
		taskViews[i].Order = i + 1
	}
	return model.TemplateTasks{
		Tasks: taskViews,
	}, err
}

func (s *Service) GetCreateTask() (model.TemplateEditTask, error) {
	dockerfiles := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&dockerfiles)
	return model.TemplateEditTask{
		Dockerfiles: dockerfiles,
	}, err
}

func (s *Service) GetEditTask(id int) (model.TemplateEditTask, error) {
	result := model.TemplateEditTask{
		IsEdit: true,
	}
	dockerfiles := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&dockerfiles)
	result.Dockerfiles = dockerfiles
	if err != nil {
		return result, fmt.Errorf("err get dockerfile templates: %w", err)
	}
	var task model.TaskDB
	err = s.DB.GetEntityByField("id", strconv.Itoa(id), &task)
	if err != nil {
		return result, fmt.Errorf("err get task from db: %w", err)
	}
	result.Fill(task)
	return result, nil
}

func (s *Service) GetAttemptTask(id int) (model.TemplateUploadAttempt, error) {
	result := model.TemplateUploadAttempt{}
	var task model.TaskDB
	err := s.DB.GetEntityByField("id", strconv.Itoa(id), &task)
	if err != nil {
		return result, fmt.Errorf("err get task from db: %w", err)
	}
	result.Fill(task)
	return result, nil
}

func (s *Service) CreateTask(req model.ChangeTaskRequest) error {
	task := model.TaskDB{}
	task.Fill(req)
	err := s.DB.CreateEntity(&task)
	if err != nil {
		return fmt.Errorf("err create task: %w", err)
	}
	return nil
}

func (s *Service) EditTask(req model.ChangeTaskRequest) error {
	var task model.TaskDB
	err := s.DB.GetEntityByField("id", strconv.Itoa(req.ID), &task)
	if err != nil {
		return fmt.Errorf("err get task from db: %w", err)
	}
	task.Fill(req)
	err = s.DB.UpdateEntity(&task)
	if err != nil {
		return fmt.Errorf("err update task: %w", err)
	}
	return nil
}

func (s *Service) DeleteTask(id int) error {
	err := s.DB.DeleteEntityByField("id", strconv.Itoa(id), model.TaskDB{})
	if err != nil {
		return fmt.Errorf("err delete task in db: %w", err)
	}
	return nil
}

func (s *Service) GetAttempts(userID string) (model.TemplateAttempts, error) {
	attempts, err := s.DB.GetAttemptsByStudentID(userID)
	for i := range attempts {
		attempts[i].Order = i + 1
		attempts[i].CreatedAtFormat = attempts[i].CreatedAt.Format("2006-01-02 15:04:05")
	}
	return model.TemplateAttempts{
		Attempts: attempts,
	}, err
}

func (s *Service) GetDockerfileTemplates() ([]model.DockerfileTemplateDB, error) {
	entities := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&entities)
	return entities, err
}

func (s *Service) GetAttempt(attemptID string, isAuthorized bool, role string) (model.TemplateAttempt, error) {
	attempt, err := s.DB.GetAttemptByID(attemptID)
	attempt.CreatedAtFormat = attempt.CreatedAt.Format("2006-01-02 15:04:05")
	return model.TemplateAttempt{
		Attempt:      attempt,
		IsAuthorized: isAuthorized,
		Role:         role,
	}, err
}
