package db

import (
	"github.com/nightlord189/ulp/internal/model"
)

func (d *Manager) GetTasksByCreatorID(creatorID int) ([]model.TaskDB, error) {
	result := make([]model.TaskDB, 0)
	err := d.DB.Where("creator_id = ?", creatorID).Order("id DESC").Find(&result).Error
	return result, err
}
