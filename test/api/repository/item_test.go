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

func TestItemRepository_FindByProjectIdAndItemId_Success(t *testing.T) {
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

	item, _ := itemRepository.FindByProjectIdAndItemId(123, 345)
	cleanItemDateAndTime(item)

	expectedItem := models.Item{
		ID:          345,
		ProjectID:   123,
		Name:        "itemName",
		Description: "itemDescription",
	}

	assert.Equal(t, expectedItem, *item)
}

func TestItemRepository_FindByProjectIdAndItemId_RecordNotFound(t *testing.T) {
	_, itemRepository, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	item, err := itemRepository.FindByProjectIdAndItemId(123, 345)

	assert.Equal(t, &models.Item{}, item)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestItemRepository_Update_Success(t *testing.T) {
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

	itemRepository.Update(123, 345, map[string]interface{}{
		"name": "newName",
	})

	item, _ := itemRepository.FindByProjectIdAndItemId(123, 345)
	cleanItemDateAndTime(item)

	expectedItem := models.Item{
		ID:          345,
		ProjectID:   123,
		Name:        "newName",
		Description: "itemDescription",
	}

	assert.Equal(t, expectedItem, *item)
}

func TestItemRepository_Update_RecordNotFound(t *testing.T) {
	_, itemRepository, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	err := itemRepository.Update(123, 345, map[string]interface{}{
		"name": "newName",
	})

	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestItemRepository_Delete_Success(t *testing.T) {
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

	errDelete := itemRepository.Delete(123, 345)
	item, errFind := itemRepository.FindByProjectIdAndItemId(123, 345)
	cleanItemDateAndTime(item)

	assert.Nil(t, errDelete)
	assert.Equal(t, &models.Item{}, item)
	assert.True(t, errors.Is(errFind, gorm.ErrRecordNotFound))
}

func TestItemRepository_Delete_RecordNotFound(t *testing.T) {
	_, itemRepository, teardownFunc := testutils.SetupTestProjectAndItemRepository(t)
	defer teardownFunc(t)

	errDelete := itemRepository.Delete(123, 345)

	assert.True(t, errors.Is(errDelete, gorm.ErrRecordNotFound))
}

func cleanItemDateAndTime(item *models.Item) {
	item.Model.CreatedAt = time.Time{}
	item.Model.UpdatedAt = time.Time{}
	item.Model.DeletedAt = gorm.DeletedAt{}
}
