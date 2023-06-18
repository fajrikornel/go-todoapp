package testutils

import (
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"gorm.io/gorm"
	"testing"
)

func SetupTestProjectAndItemRepository(tb testing.TB) (repository.ProjectRepository, repository.ItemRepository, func(tb testing.TB)) {
	testDbConfig := config.GetTestDbConfig()

	sqlStore, err := db.GetSqlStore(&testDbConfig)
	if err != nil {
		panic("Cannot instantiate test sqlStore")
	}

	projectRepository := repository.NewProjectRepository(sqlStore)
	itemRepository := repository.NewItemRepository(sqlStore)

	teardownProjectRepositoryFunc := getTeardownFunction(sqlStore)

	return projectRepository, itemRepository, teardownProjectRepositoryFunc
}

func getTeardownFunction(sqlStore *db.SqlStore) func(tb testing.TB) {
	var projectIds []uint
	itemIds := make(map[uint]uint)

	sqlStore.Db.Callback().Create().After("gorm:create").Register("rememberCreatedObjects", func(db *gorm.DB) {
		switch db.Statement.Schema.Table {
		case "projects":
			projectIds = append(projectIds, db.Statement.Dest.(*models.Project).ID)
		case "items":
			item := db.Statement.Dest.(*models.Item)
			itemIds[item.ID] = item.ProjectID
		}
	})

	teardownFunc := func(tb testing.TB) {
		for itemId, projectId := range itemIds {
			sqlStore.Db.Unscoped().Delete(&models.Item{ID: itemId, ProjectID: projectId})
		}

		for _, projectId := range projectIds {
			sqlStore.Db.Unscoped().Delete(&models.Project{ID: projectId})
		}
	}
	return teardownFunc
}
