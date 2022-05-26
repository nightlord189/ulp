package db

import (
	"errors"
	"github.com/nightlord189/ulp/internal/model"
)

const getAttemptsSQL = "select a.*, t.name as t_name " +
	"from attempts a left join tasks t on t.id = a.task_id " +
	"where a.creator_id = ? order by a.id desc"

const getAttemptSQL = "select a.*, t.name as t_name, t.description as t_description, t.type as t_type, u.username as u_username " +
	"from attempts a " +
	"left join tasks t on t.id = a.task_id " +
	"left join users u on u.id = a.creator_id " +
	"where a.id = ? order by a.id desc"

func (d *Manager) GetAttemptsByStudentID(studentID int) ([]model.AttemptView, error) {
	entities := make([]model.AttemptView, 0)
	err := d.DB.Raw(getAttemptsSQL, studentID).Scan(&entities).Error
	return entities, err
}

func (d *Manager) GetAttemptByID(id string) (model.AttemptView, error) {
	entities := make([]model.AttemptView, 0, 1)
	err := d.DB.Raw(getAttemptSQL, id).Scan(&entities).Error
	if len(entities) < 1 {
		return model.AttemptView{}, errors.New("not found")
	}
	return entities[0], err
}
