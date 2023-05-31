package models

import (
	"gorm.io/gorm"
)

type Item struct {
	ID          uint   `json:"item_id"`
	ProjectID   uint   `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	gorm.Model  `json:"-"`
}
