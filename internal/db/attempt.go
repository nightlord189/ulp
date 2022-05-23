package db

import (
	"github.com/nightlord189/ulp/internal/model"
)

const getAttemptsSQL = "select a.*, t.name as t_name " +
	"from attempts a left join tasks t on t.id = a.task_id " +
	"where a.creator_id = ? order by a.id desc"

func (d *Manager) GetAttemptsByStudentID(studentID string) ([]model.AttemptView, error) {
	entities := make([]model.AttemptView, 0)
	err := d.DB.Raw(getAttemptsSQL, studentID).Scan(&entities).Error
	return entities, err
}
