package model

import "time"

type AttemptDB struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskID      int       `json:"taskId" gorm:"column:task_id"`
	State       string    `json:"state" gorm:"column:state"`
	Log         string    `json:"log" gorm:"column:log"`
	RunningTime int       `json:"runningTime" gorm:"column:running_time"`
	CreatorID   int       `json:"creatorID" gorm:"column:creator_id"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (AttemptDB) TableName() string {
	return "attempts"
}

type AttemptRequest struct {
	TaskID    string `json:"taskId" binding:"required"`
	CreatorID int    `json:"creatorID" binding:"required"`
}
