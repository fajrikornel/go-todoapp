package db

import (
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlStore struct {
	Db *gorm.DB
}

func GetSqlStore(conf *config.Config) (*SqlStore, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		conf.DbHost,
		conf.DbUsername,
		conf.DbPassword,
		conf.DbName,
		conf.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &SqlStore{}, err
	}

	return &SqlStore{Db: db}, nil
}

func (s *SqlStore) DoMigrations() error {
	return s.Db.AutoMigrate(&models.Project{}, &models.Item{})
}
