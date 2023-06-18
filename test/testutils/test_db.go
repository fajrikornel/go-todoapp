package testutils

import (
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"gorm.io/gorm"
	"testing"
)

func SetupTestProjectRepository(tb testing.TB) (repository.ProjectRepository, func(tb testing.TB)) {
	testDbConfig := config.GetTestDbConfig()

	sqlStore, err := db.GetSqlStore(&testDbConfig)
	if err != nil {
		panic("Cannot instantiate test sqlStore")
	}

	projectRepository := repository.NewProjectRepository(sqlStore)

	teardownProjectRepositoryFunc := getTeardownFunction(sqlStore)

	return projectRepository, teardownProjectRepositoryFunc
}

func SetupTestItemRepository(tb testing.TB) (repository.ItemRepository, func(tb testing.TB)) {
	testDbConfig := config.GetTestDbConfig()

	sqlStore, err := db.GetSqlStore(&testDbConfig)
	if err != nil {
		panic("Cannot instantiate test sqlStore")
	}

	itemRepository := repository.NewItemRepository(sqlStore)

	teardownItemRepositoryFunc := getTeardownFunction(sqlStore)

	return itemRepository, teardownItemRepositoryFunc
}

func getTeardownFunction(sqlStore *db.SqlStore) func(tb testing.TB) {
	var projectIds []uint
	var itemIds []uint

	sqlStore.Db.Callback().Create().After("gorm:create").Register("rememberCreatedObjects", func(db *gorm.DB) {
		switch db.Statement.Schema.Table {
		case "projects":
			projectIds = append(projectIds, db.Statement.Dest.(*models.Project).ID)
		case "items":
			itemIds = append(itemIds, db.Statement.Dest.(*models.Item).ID)
		}
	})

	teardownProjectRepositoryFunc := func(tb testing.TB) {
		for _, projectId := range projectIds {
			sqlStore.Db.Unscoped().Delete(&models.Project{ID: projectId})
		}

		for _, itemId := range itemIds {
			sqlStore.Db.Unscoped().Delete(&models.Project{ID: itemId})
		}
	}
	return teardownProjectRepositoryFunc
}
