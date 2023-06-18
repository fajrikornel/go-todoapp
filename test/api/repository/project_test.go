package repository

import (
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/test/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestProjectRepository_FindById_Success(t *testing.T) {
	projectRepository, itemRepository, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	projectRepository.Create(&models.Project{
		ID:          123,
		Name:        "projectName",
		Description: "projectDescription",
	})

	itemRepository.Create(&models.Item{
		ID:          345,
		ProjectID:   123,
		Name:        "itemName",
		Description: "itemDescription",
	})

	project, _ := projectRepository.FindById(123)
	cleanProjectDateAndTime(project)

	expectedProject := models.Project{
		ID:          123,
		Name:        "projectName",
		Description: "projectDescription",
		Items: []models.Item{
			{
				ID:          345,
				ProjectID:   123,
				Name:        "itemName",
				Description: "itemDescription",
			},
		},
	}

	assert.Equal(t, expectedProject, *project)
}

func TestProjectRepository_FindById_RecordNotFound(t *testing.T) {
	projectRepository, _, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	project, err := projectRepository.FindById(123)

	assert.Equal(t, &models.Project{}, project)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestProjectRepository_Update_Success(t *testing.T) {
	projectRepository, _, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	projectRepository.Create(&models.Project{
		ID:          123,
		Name:        "projectName",
		Description: "projectDescription",
	})

	projectRepository.Update(123, map[string]interface{}{
		"name": "newName",
	})

	project, _ := projectRepository.FindById(123)
	cleanProjectDateAndTime(project)

	expectedProject := models.Project{
		ID:          123,
		Name:        "newName",
		Description: "projectDescription",
		Items:       []models.Item{},
	}

	assert.Equal(t, expectedProject, *project)
}

func TestProjectRepository_Update_RecordNotFound(t *testing.T) {
	projectRepository, _, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	err := projectRepository.Update(123, map[string]interface{}{
		"name": "newName",
	})

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestProjectRepository_Delete_Success(t *testing.T) {
	projectRepository, _, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	projectRepository.Create(&models.Project{
		ID:          123,
		Name:        "projectName",
		Description: "projectDescription",
	})

	errDelete := projectRepository.Delete(123)
	project, errFind := projectRepository.FindById(123)
	cleanProjectDateAndTime(project)

	assert.Nil(t, errDelete)
	assert.Equal(t, &models.Project{}, project)
	assert.True(t, errors.Is(errFind, gorm.ErrRecordNotFound))
}

func TestProjectRepository_Delete_RecordNotFound(t *testing.T) {
	projectRepository, _, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	errDelete := projectRepository.Delete(123)

	assert.True(t, errors.Is(errDelete, gorm.ErrRecordNotFound))
}

func cleanProjectDateAndTime(project *models.Project) {
	project.Model.CreatedAt = time.Time{}
	project.Model.UpdatedAt = time.Time{}
	project.Model.DeletedAt = gorm.DeletedAt{}

	for i, _ := range project.Items {
		item := &project.Items[i]
		item.Model.CreatedAt = time.Time{}
		item.Model.UpdatedAt = time.Time{}
		item.Model.DeletedAt = gorm.DeletedAt{}
	}
}
