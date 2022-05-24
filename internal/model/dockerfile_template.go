package model

type DockerfileTemplateDB struct {
	ID      int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Key     string `json:"key" gorm:"column:key;unique"`
	Name    string `json:"name" gorm:"column:name;not null"`
	Content string `json:"content" gorm:"column:content"`
}

func (DockerfileTemplateDB) TableName() string {
	return "dockerfile_templates"
}
