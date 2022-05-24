package db

import (
	"fmt"
)

func (d *Manager) DeleteEntityByField(field, value string, entityModel interface{}) error {
	return d.DB.Where(fmt.Sprintf("%s = ?", field), value).Delete(entityModel).Error
}

func (d *Manager) DeleteEntityByFields(fields map[string]interface{}, entityModel interface{}) error {
	return d.DB.Where(fields).Delete(entityModel).Error
}

func (d *Manager) UpdateEntity(entity interface{}) error {
	return d.DB.Save(entity).Error
}

func (d *Manager) GetEntityByField(field, value string, entity interface{}) error {
	return d.DB.Where(fmt.Sprintf("%s = ?", field), value).First(entity).Error
}

func (d *Manager) GetEntityByFields(fields map[string]interface{}, entity interface{}) error {
	return d.DB.Where(fields).First(entity).Error
}

func (d *Manager) GetAllEntities(entities interface{}) error {
	return d.DB.Find(entities).Error
}

func (d *Manager) GetEntitiesByField(field, value string, entity interface{}) error {
	return d.DB.Where(fmt.Sprintf("%s = ?", field), value).Find(entity).Error
}

func (d *Manager) GetEntitiesByFieldQuery(field, value string, entity interface{}, limit, offset int) error {
	query := d.DB.Limit(limit).Offset(offset).Where(fmt.Sprintf("%s = ?", field), value)
	if limit > 0 || offset > 0 {
		query = query.Order("id")
	}
	return query.Find(entity).Error
}

func (d *Manager) GetEntitiesByFields(fields map[string]interface{}, entities interface{}) error {
	return d.DB.Where(fields).Find(entities).Error
}

func (d *Manager) GetEntitiesByFieldsQuery(fields map[string]interface{}, entities interface{}, limit, offset int) error {
	query := d.DB.Where(fields)
	if limit > 0 || offset > 0 {
		query = query.Order("id")
	}
	return query.Find(entities).Error
}

func (d *Manager) CreateEntities(entities interface{}) error {
	return d.DB.Create(entities).Error
}

func (d *Manager) CreateEntity(entity interface{}) error {
	return d.DB.Create(entity).Error
}

func (d *Manager) GetSequenceValue(sequence string) (string, error) {
	result := ""
	err := d.DB.Raw("SELECT nextval(?)", sequence).Scan(&result).Error
	return result, err
}
