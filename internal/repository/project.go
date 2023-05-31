package repository

import (
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindById(id int) (*models.Project, error)
}

type projectRepository struct {
	sqlStore *db.SqlStore
}

func NewProjectRepository(sqlStore *db.SqlStore) ProjectRepository {
	return &projectRepository{
		sqlStore: sqlStore,
	}
}

func (p *projectRepository) Create(project *models.Project) error {
	tx := p.sqlStore.Db.Create(project)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (p *projectRepository) FindById(id int) (*models.Project, error) {
	project := &models.Project{}
	tx := p.sqlStore.Db.Preload("Items").First(project, id)
	if tx.Error != nil {
		return project, tx.Error
	}

	return project, nil
}
