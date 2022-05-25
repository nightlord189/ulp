package model

import "time"

const (
	AttemptStateSuccess AttemptState = "success"
	AttemptStateFail    AttemptState = "fail"
)

type AttemptState string

type AttemptDB struct {
	ID          int          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskID      int          `json:"taskId" gorm:"column:task_id"`
	State       AttemptState `json:"state" gorm:"column:state"`
	Log         string       `json:"log" gorm:"column:log"`
	RunningTime int          `json:"runningTime" gorm:"column:running_time"`
	CreatorID   int          `json:"creatorID" gorm:"column:creator_id"`
	CreatedAt   time.Time    `json:"createdAt"`
}

func (AttemptDB) TableName() string {
	return "attempts"
}

type AttemptView struct {
	Order           int       `gorm:"-"`
	ID              int       `gorm:"column:id"`
	TaskName        string    `gorm:"column:t_name"`
	TaskDescription string    `gorm:"column:t_description"`
	TaskType        string    `gorm:"column:t_type"`
	TaskID          int       `gorm:"column:task_id"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	CreatedAtFormat string    `gorm:"-"`
	State           string    `gorm:"column:state"`
	CreatorUsername string    `gorm:"column:u_username"`
	RunningTime     int       `gorm:"column:running_time"`
	Log             string    `gorm:"column:log"`
}

type AttemptRequest struct {
	TaskID    string `json:"taskId" binding:"required"`
	CreatorID int    `json:"creatorID" binding:"required"`
}
