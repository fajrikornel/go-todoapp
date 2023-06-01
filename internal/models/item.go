package models

import (
	"gorm.io/gorm"
)

type Item struct {
	ID          uint   `json:"item_id" gorm:"primaryKey;autoIncrement:true"`
	ProjectID   uint   `json:"-" gorm:"primaryKey;autoIncrement:false"`
	Name        string `json:"name"`
	Description string `json:"description"`
	gorm.Model  `json:"-"`
}
