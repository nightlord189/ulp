package service

import (
	"errors"
	"fmt"
	"github.com/nightlord189/ulp/internal/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

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

func (s *Service) GetTasks(userID int, userInfo model.TemplateUserInfo) (model.TemplateTasks, error) {
	tasks, err := s.DB.GetTasksByCreatorID(userID)
	taskViews := make([]model.TaskView, len(tasks))
	for i := range tasks {
		taskViews[i] = tasks[i].ToView()
		taskViews[i].Order = i + 1
	}
	return model.TemplateTasks{
		Tasks:    taskViews,
		UserInfo: userInfo,
	}, err
}

func (s *Service) GetCreateTask(userID int, userInfo model.TemplateUserInfo) (model.TemplateEditTask, error) {
	dockerfiles := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&dockerfiles)
	return model.TemplateEditTask{
		UserID:      userID,
		Dockerfiles: dockerfiles,
		UserInfo:    userInfo,
	}, err
}

func (s *Service) GetEditTask(taskID, userID int, userInfo model.TemplateUserInfo) (model.TemplateEditTask, error) {
	result := model.TemplateEditTask{
		UserID:   userID,
		IsEdit:   true,
		UserInfo: userInfo,
	}
	dockerfiles := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&dockerfiles)
	result.Dockerfiles = dockerfiles
	if err != nil {
		return result, fmt.Errorf("err get dockerfile templates: %w", err)
	}
	var task model.TaskDB
	err = s.DB.GetEntityByField("id", strconv.Itoa(taskID), &task)
	if err != nil {
		return result, fmt.Errorf("err get task from db: %w", err)
	}
	result.Fill(task)
	return result, nil
}

func (s *Service) GetAttemptTask(taskID, userID int, userInfo model.TemplateUserInfo) (model.TemplateUploadAttempt, error) {
	result := model.TemplateUploadAttempt{
		UserID:   userID,
		UserInfo: userInfo,
	}
	var task model.TaskDB
	err := s.DB.GetEntityByField("id", strconv.Itoa(taskID), &task)
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
	if task.CreatorID != req.CreatorID {
		return errors.New("задание может редактировать только его создатель")
	}
	task.Fill(req)
	err = s.DB.UpdateEntity(&task)
	if err != nil {
		return fmt.Errorf("err update task: %w", err)
	}
	return nil
}

func (s *Service) DeleteTask(taskID, userID int) error {
	var task model.TaskDB
	err := s.DB.GetEntityByField("id", strconv.Itoa(taskID), &task)
	if err != nil {
		return fmt.Errorf("err get task from db: %w", err)
	}
	if task.CreatorID != userID {
		return errors.New("задание может удалять только его создатель")
	}
	err = s.DB.DeleteEntityByField("id", strconv.Itoa(taskID), model.TaskDB{})
	if err != nil {
		return fmt.Errorf("err delete task in db: %w", err)
	}
	return nil
}

func (s *Service) GetAttempts(userID int, userInfo model.TemplateUserInfo) (model.TemplateAttempts, error) {
	attempts, err := s.DB.GetAttemptsByStudentID(userID)
	for i := range attempts {
		attempts[i].Order = i + 1
		attempts[i].CreatedAtFormat = attempts[i].CreatedAt.Format("2006-01-02 15:04:05")
	}
	return model.TemplateAttempts{
		Attempts: attempts,
		UserInfo: userInfo,
	}, err
}

func (s *Service) GetDockerfileTemplates() ([]model.DockerfileTemplateDB, error) {
	entities := make([]model.DockerfileTemplateDB, 0)
	err := s.DB.GetAllEntities(&entities)
	return entities, err
}

func (s *Service) GetAttempt(attemptID string, userInfo model.TemplateUserInfo) (model.TemplateAttempt, error) {
	attempt, err := s.DB.GetAttemptByID(attemptID)
	attempt.CreatedAtFormat = attempt.CreatedAt.Format("2006-01-02 15:04:05")
	return model.TemplateAttempt{
		Attempt:  attempt,
		UserInfo: userInfo,
	}, err
}
