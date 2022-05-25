package model

import (
	"errors"
	"time"
)

const (
	TaskTypeConsole TaskType = "console"
	TaskTypeWebAPI           = "web_api"
	TaskTypeHTML             = "html"
)

type TaskType string

type TaskDB struct {
	ID               int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name             string    `json:"name" gorm:"column:name;unique"`
	Description      string    `json:"description" gorm:"column:description"`
	Type             TaskType  `json:"type" gorm:"column:type"`
	Dockerfile       string    `json:"dockerfile" gorm:"column:dockerfile"`
	CreatorID        int       `json:"creatorID" gorm:"column:creator_id"`
	TestcaseType     string    `json:"testCaseType" gorm:"column:testcase_type"`
	TestcaseURL      string    `json:"testCaseURL" gorm:"column:testcase_url"`
	TestcaseExpected string    `json:"testCaseExpected" gorm:"column:testcase_expected"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (TaskDB) TableName() string {
	return "tasks"
}

type TaskView struct {
	Order       int
	ID          int
	Name        string
	Description string
	Type        string
	CreatedAt   string
	UpdatedAt   string
}

func (t *TaskDB) ToView() TaskView {
	taskType := ""
	switch t.Type {
	case TaskTypeConsole:
		taskType = "console"
		break
	case TaskTypeWebAPI:
		taskType = "Web-API"
		break
	case TaskTypeHTML:
		taskType = "HTML"
		break
	}

	return TaskView{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Type:        taskType,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (t *TaskDB) Fill(req ChangeTaskRequest) {
	t.Name = req.Name
	t.Description = req.Description
	t.CreatorID = req.CreatorID
	t.Type = req.Type
	t.Dockerfile = req.Dockerfile
	t.TestcaseType = req.TestcaseType
	if t.Type == TaskTypeWebAPI || t.Type == TaskTypeHTML {
		t.TestcaseURL = req.TestcaseURL
	} else {
		t.TestcaseURL = ""
	}
	t.TestcaseExpected = req.TestcaseExpected
}

type ChangeTaskRequest struct {
	ID               int
	Name             string   `json:"name" binding:"required" form:"name"`
	Description      string   `json:"description" form:"description"`
	CreatorID        int      `json:"creatorID"`
	Type             TaskType `json:"type" binding:"required" form:"taskType"`
	Dockerfile       string   `json:"dockerfile" binding:"required" form:"dockerfile"`
	TestcaseType     string   `json:"testCaseType" binding:"required" form:"testCaseType"`
	TestcaseURL      string   `json:"testCaseURL" form:"testCaseUrl"`
	TestcaseExpected string   `json:"testCaseExpected" binding:"required" form:"testCaseExpected"`
}

func (r *ChangeTaskRequest) IsValid() error {
	if r.Name == "" {
		return errors.New("name is empty")
	}
	if r.Type == "" {
		return errors.New("type is empty")
	}
	if r.Dockerfile == "" {
		return errors.New("dockerfile is empty")
	}
	if r.TestcaseType == "" {
		return errors.New("testcase type is empty")
	}
	if r.TestcaseExpected == "" {
		return errors.New("testcase expected is empty")
	}
	if (r.TestcaseType == TaskTypeWebAPI || r.TestcaseType == TaskTypeHTML) && r.TestcaseURL == "" {
		return errors.New("testcase URL is empty")
	}
	return nil
}
