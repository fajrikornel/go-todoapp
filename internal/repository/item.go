package repository

import (
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
)

type ItemRepository interface {
	Create(item *models.Item) error
	FindByProjectIdAndItemId(projectId, itemId int) (*models.Item, error)
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
