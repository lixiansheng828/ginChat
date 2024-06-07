package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId string
	Type    string
	Desc    int
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
