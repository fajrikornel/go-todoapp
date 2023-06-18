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

func TestGetItemHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := GetItemHandler(mItemRepository)

	router := httprouter.New()
	router.GET("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName             string
		projectId            int
		itemId               int
		returnedError        error
		expectedErrorMessage string
		expectedHttpCode     int
	}{
		{
			"project does not exist in database",
			123,
			345,
			gorm.ErrRecordNotFound,
			"item_or_project_not_found",
			400,
		},
		{
			"error while calling database",
			123,
			345,
			gorm.ErrUnsupportedDriver,
			"internal_db_error",
			500,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mItemRepository.
				EXPECT().
				FindByProjectIdAndItemId(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId)).
				Return(nil, tc.returnedError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.expectedHttpCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.expectedHttpCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[GetItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[GetItemResponseData]{
				Success: false,
				Error:   tc.expectedErrorMessage,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestGetItemHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := GetItemHandler(mItemRepository)

	router := httprouter.New()
	router.GET("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName     string
		projectId    int
		itemId       int
		existingItem models.Item
	}{
		{
			"item exists",
			123,
			345,
			models.Item{
				ID:          345,
				ProjectID:   123,
				Name:        "itemName",
				Description: "itemDescription",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			returnedItem := models.Item{
				ID:          tc.existingItem.ID,
				ProjectID:   tc.existingItem.ProjectID,
				Name:        tc.existingItem.Name,
				Description: tc.existingItem.Description,
			}

			mItemRepository.
				EXPECT().
				FindByProjectIdAndItemId(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId)).
				Return(&returnedItem, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[GetItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[GetItemResponseData]{
				Success: true,
				Data: GetItemResponseData{
					ItemID:      uint(tc.itemId),
					Name:        tc.existingItem.Name,
					Description: tc.existingItem.Description,
				},
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
