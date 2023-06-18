package v1

import (
	"encoding/json"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	. "github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/models"
	mock_repository "github.com/fajrikornel/go-todoapp/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetProjectHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := GetProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.GET("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName         string
		projectId        int
		returnedError    error
		expectedHttpCode int
	}{
		{
			"project does not exist in database",
			123,
			gorm.ErrRecordNotFound,
			400,
		},
		{
			"error while calling database",
			123,
			gorm.ErrUnsupportedDriver,
			500,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mProjectRepository.
				EXPECT().
				FindById(gomock.Eq(tc.projectId)).
				Return(nil, tc.returnedError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/projects/%d", tc.projectId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.expectedHttpCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.expectedHttpCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[GetProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[GetProjectResponseData]{
				Success: false,
				Error:   tc.returnedError.Error(),
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestGetProjectHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := GetProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.GET("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName        string
		projectId       int
		existingProject models.Project
		existingItems   []models.Item
	}{
		{
			"project exists but with no items",
			123,
			models.Project{
				ID:          123,
				Name:        "projectName",
				Description: "projectDescription",
			},
			[]models.Item{},
		},
		{
			"project exists with one item",
			123,
			models.Project{
				ID:          123,
				Name:        "projectName",
				Description: "projectDescription",
			},
			[]models.Item{
				{
					ID:          345,
					ProjectID:   123,
					Name:        "itemName",
					Description: "itemDescription",
				},
			},
		},
		{
			"project exists with multiple items",
			123,
			models.Project{
				ID:          123,
				Name:        "projectName",
				Description: "projectDescription",
			},
			[]models.Item{
				{
					ID:          345,
					ProjectID:   123,
					Name:        "itemName",
					Description: "itemDescription",
				},
				{
					ID:          346,
					ProjectID:   123,
					Name:        "itemName_2",
					Description: "itemDescription_2",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			returnedProject := models.Project{
				ID:          tc.existingProject.ID,
				Name:        tc.existingProject.Name,
				Description: tc.existingProject.Description,
				Items:       tc.existingItems,
			}

			mProjectRepository.
				EXPECT().
				FindById(gomock.Eq(tc.projectId)).
				Return(&returnedProject, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/projects/%d", tc.projectId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 400, rr.Code)
			}

			var actualResponse utils.GenericResponse[GetProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[GetProjectResponseData]{
				Success: true,
				Data: GetProjectResponseData{
					ProjectID:   uint(tc.projectId),
					Name:        tc.existingProject.Name,
					Description: tc.existingProject.Description,
					Items:       convertExistingItemsToItemResponseStruct(tc.existingItems),
				},
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func convertExistingItemsToItemResponseStruct(items []models.Item) []ItemResponse {
	var itemResponses []ItemResponse
	for _, i := range items {
		itemResponse := ItemResponse{
			ItemID: i.ID,
			Name:   i.Name,
		}
		itemResponses = append(itemResponses, itemResponse)
	}
	return itemResponses
}
