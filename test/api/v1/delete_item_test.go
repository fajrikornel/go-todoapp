package v1

import (
	"encoding/json"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	. "github.com/fajrikornel/go-todoapp/internal/api/v1"
	mock_repository "github.com/fajrikornel/go-todoapp/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDeleteItemHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := DeleteItemHandler(mItemRepository)

	router := httprouter.New()
	router.DELETE("/v1/projects/:projectId/:itemId", handleFunc)

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
			404,
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
				Delete(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId)).
				Return(tc.returnedError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.expectedHttpCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.expectedHttpCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[DeleteItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[DeleteItemResponseData]{
				Success: false,
				Error:   tc.expectedErrorMessage,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestDeleteItemHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := DeleteItemHandler(mItemRepository)

	router := httprouter.New()
	router.DELETE("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName  string
		projectId int
		itemId    int
	}{
		{
			"success deleting project",
			123,
			345,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mItemRepository.
				EXPECT().
				Delete(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId)).
				Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[DeleteItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[DeleteItemResponseData]{
				Success: true,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
