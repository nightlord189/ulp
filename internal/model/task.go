package model

import (
	"errors"
	"time"
)

const (
	TaskTypeConsole      TaskType     = "console"
	TaskTypeWebAPI       TaskType     = "web_api"
	TaskTypeHTML         TaskType     = "html"
	TestCaseTypeContains TestCaseType = "contains"
	TestCaseTypeEqual    TestCaseType = "equal"
)

type TaskType string

type TestCaseType string

type TaskDB struct {
	ID               int          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name             string       `json:"name" gorm:"column:name;unique"`
	Description      string       `json:"description" gorm:"column:description"`
	Type             TaskType     `json:"type" gorm:"column:type"`
	Dockerfile       string       `json:"dockerfile" gorm:"column:dockerfile"`
	CreatorID        int          `json:"creatorID" gorm:"column:creator_id"`
	TestcaseType     TestCaseType `json:"testCaseType" gorm:"column:testcase_type"`
	TestcaseURL      string       `json:"testCaseURL" gorm:"column:testcase_url"`
	TestcaseInput    string       `json:"testCaseInput" gorm:"column:testcase_input"`
	TestcaseExpected string       `json:"testCaseExpected" gorm:"column:testcase_expected"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	CreatedAt        time.Time    `json:"createdAt"`
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
	t.TestcaseInput = req.TestcaseInput
	t.TestcaseExpected = req.TestcaseExpected
}

type ChangeTaskRequest struct {
	ID               int
	Name             string       `json:"name" binding:"required" form:"name"`
	Description      string       `json:"description" form:"description"`
	CreatorID        int          `json:"userID" form:"userID"`
	Type             TaskType     `json:"type" binding:"required" form:"taskType"`
	Dockerfile       string       `json:"dockerfile" binding:"required" form:"dockerfile"`
	TestcaseType     TestCaseType `json:"testCaseType" binding:"required" form:"testCaseType"`
	TestcaseURL      string       `json:"testCaseURL" form:"testCaseUrl"`
	TestcaseInput    string       `json:"testCaseInput" form:"testcaseInput"`
	TestcaseExpected string       `json:"testCaseExpected" binding:"required" form:"testCaseExpected"`
}

func (r *ChangeTaskRequest) IsValid() error {
	if r.Name == "" {
		return errors.New("name is empty")
	}
	if r.Dockerfile == "" {
		return errors.New("dockerfile is empty")
	}
	if r.TestcaseExpected == "" {
		return errors.New("testcase expected is empty")
	}
	if r.Type != TaskTypeConsole && r.Type != TaskTypeWebAPI && r.Type != TaskTypeHTML {
		return errors.New("invalid task type format")
	}
	if r.TestcaseType != TestCaseTypeContains && r.TestcaseType != TestCaseTypeEqual {
		return errors.New("invalid testcase type format")
	}
	if (r.Type == TaskTypeWebAPI || r.Type == TaskTypeHTML) && r.TestcaseURL == "" {
		return errors.New("testcase URL is empty")
	}
	if r.Type == TaskTypeConsole && r.TestcaseInput == "" {
		return errors.New("testcase input is empty")
	}
	return nil
}
