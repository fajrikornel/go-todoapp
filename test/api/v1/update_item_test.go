package v1

import (
	"encoding/json"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	. "github.com/fajrikornel/go-todoapp/internal/api/v1"
	mock_repository "github.com/fajrikornel/go-todoapp/test/mocks/repository"
	"github.com/fajrikornel/go-todoapp/test/testutils"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUpdateItemHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := UpdateItemHandler(mItemRepository)

	router := httprouter.New()
	router.PATCH("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName             string
		projectId            int
		itemId               int
		name                 *string
		description          *string
		returnedError        error
		expectedErrorMessage string
		expectedHttpCode     int
	}{
		{
			"project does not exist in database",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			testutils.CreatePointerOfString("description"),
			gorm.ErrRecordNotFound,
			"item_or_project_not_found",
			400,
		},
		{
			"error while calling database",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			testutils.CreatePointerOfString("description"),
			gorm.ErrUnsupportedDriver,
			"internal_db_error",
			500,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := UpdateItemRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			fields := testutils.ConstructUpdateFieldsMap(tc.name, tc.description)
			mItemRepository.
				EXPECT().
				Update(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId), gomock.Eq(fields)).
				Return(tc.returnedError)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), testutils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.expectedHttpCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.expectedHttpCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[UpdateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[UpdateItemResponseData]{
				Success: false,
				Error:   tc.expectedErrorMessage,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestUpdateItemHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := UpdateItemHandler(mItemRepository)

	router := httprouter.New()
	router.PATCH("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName    string
		projectId   int
		itemId      int
		name        *string
		description *string
		error       string
	}{
		{
			"both name and description does not exist",
			123,
			345,
			nil,
			nil,
			"name_and_description_empty",
		},
		{
			"both name and description are empty strings",
			123,
			345,
			testutils.CreatePointerOfString(""),
			testutils.CreatePointerOfString(""),
			"name_and_description_empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := UpdateItemRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mItemRepository.
				EXPECT().
				Update(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId), gomock.Any()).
				Times(0)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), testutils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 400 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 400, rr.Code)
			}

			var actualResponse utils.GenericResponse[UpdateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[UpdateItemResponseData]{
				Success: false,
				Error:   tc.error,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestUpdateItemHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := UpdateItemHandler(mItemRepository)

	router := httprouter.New()
	router.PATCH("/v1/projects/:projectId/:itemId", handleFunc)

	testCases := []struct {
		testName             string
		projectId            int
		itemId               int
		name                 *string
		description          *string
		expectedFieldUpdates map[string]interface{}
	}{
		{
			"update both name and description",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			testutils.CreatePointerOfString("description"),
			map[string]interface{}{
				"name":        testutils.CreatePointerOfString("name"),
				"description": testutils.CreatePointerOfString("description"),
			},
		},
		{
			"update only name, description does not exist",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			nil,
			map[string]interface{}{
				"name": testutils.CreatePointerOfString("name"),
			},
		},
		{
			"update only name, description is an empty string",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			testutils.CreatePointerOfString(""),
			map[string]interface{}{
				"name": testutils.CreatePointerOfString("name"),
			},
		},
		{
			"update only description, name does not exist",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			nil,
			map[string]interface{}{
				"name": testutils.CreatePointerOfString("name"),
			},
		},
		{
			"update only description, name is an empty string",
			123,
			345,
			testutils.CreatePointerOfString("name"),
			testutils.CreatePointerOfString(""),
			map[string]interface{}{
				"name": testutils.CreatePointerOfString("name"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := UpdateItemRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mItemRepository.
				EXPECT().
				Update(gomock.Eq(tc.projectId), gomock.Eq(tc.itemId), gomock.Eq(tc.expectedFieldUpdates)).
				Return(nil)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/v1/projects/%d/%d", tc.projectId, tc.itemId), testutils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[UpdateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[UpdateItemResponseData]{
				Success: true,
				Data:    UpdateItemResponseData{},
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
