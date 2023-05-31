package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Items       []Item `json:"items"`
}
