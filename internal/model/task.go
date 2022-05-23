package model

import "time"

const (
	TaskTypeConsole TaskType = "console"
	TaskTypeWebAPI           = "web_api"
	TaskTypeHTML             = "html"
)

type TaskType string

type TaskDB struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"column:name"`
	Description   string    `json:"description" gorm:"column:description"`
	Type          TaskType  `json:"type" gorm:"column:type"`
	Dockerfile    string    `json:"dockerfile" gorm:"column:dockerfile"`
	CreatorID     int       `json:"creatorID" gorm:"column:creator_id"`
	TestcaseType  string    `json:"testCaseType" gorm:"column:testcase_type"`
	TestcaseURL   string    `json:"testCaseURL" gorm:"column:testcase_url"`
	TestcaseValue string    `json:"testCaseValue" gorm:"column:testcase_value"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
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
	}
}

type ChangeTaskRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	CreatorID     int      `json:"creatorID" binding:"required"`
	Type          TaskType `json:"type" binding:"required"`
	Dockerfile    string   `json:"dockerfile" binding:"required"`
	TestcaseType  string   `json:"testCaseType" binding:"required"`
	TestcaseURL   string   `json:"testCaseURL"`
	TestcaseValue string   `json:"testCaseValue" binding:"required"`
}
