package model

type Category struct {
	Base
	Name        string `json:"name" gorm:"column:name;"`
	Description string `json:"description" gorm:"column:description;omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
