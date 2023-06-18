package repository

import (
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *models.Item) error
	FindByProjectIdAndItemId(projectId, itemId int) (*models.Item, error)
	Update(projectId int, itemId int, fields map[string]interface{}) error
	Delete(projectId int, itemId int) error
}

type itemRepository struct {
	sqlStore *db.SqlStore
}

func NewItemRepository(sqlStore *db.SqlStore) ItemRepository {
	return &itemRepository{
		sqlStore: sqlStore,
	}
}

func (i *itemRepository) Create(item *models.Item) error {
	tx := i.sqlStore.Db.Create(item)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (i *itemRepository) FindByProjectIdAndItemId(projectId, itemId int) (*models.Item, error) {
	var item models.Item
	tx := i.sqlStore.Db.Where(&models.Item{ProjectID: uint(projectId), ID: uint(itemId)}).First(&item)
	if tx.Error != nil {
		return &item, tx.Error
	}

	return &item, nil
}

func (i *itemRepository) Update(projectId int, itemId int, fields map[string]interface{}) error {
	tx := i.sqlStore.Db.Model(&models.Item{ID: uint(itemId), ProjectID: uint(projectId)}).Updates(fields)
	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return tx.Error
}

func (i *itemRepository) Delete(projectId int, itemId int) error {
	tx := i.sqlStore.Db.Delete(&models.Item{ID: uint(itemId), ProjectID: uint(projectId)})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
