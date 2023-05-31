package models

import (
	"gorm.io/gorm"
)

type Project struct {
	ID          uint   `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Items       []Item `json:"items"`
	gorm.Model  `json:"-"`
}
