package models

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   uint   `json:"-"`
}
