package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	. "github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/models"
	mock_repository "github.com/fajrikornel/go-todoapp/test/mocks/repository"
	"github.com/fajrikornel/go-todoapp/test/test_utils"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateItemHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := CreateItemHandler(mItemRepository)

	router := httprouter.New()
	router.POST("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
		projectId   int
		error       string
	}{
		{
			"name and description does not exist",
			nil,
			nil,
			123,
			"name_or_description_empty",
		},
		{
			"name exists but description does not exist",
			test_utils.CreatePointerOfString("name"),
			nil,
			123,
			"name_or_description_empty",
		},
		{
			"description exists but name does not exist",
			nil,
			test_utils.CreatePointerOfString("description"),
			123,
			"name_or_description_empty",
		},
		{
			"name and description are empty strings",
			test_utils.CreatePointerOfString(""),
			test_utils.CreatePointerOfString(""),
			123,
			"name_or_description_empty",
		},
		{
			"name exists but description is an empty string",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString(""),
			123,
			"name_or_description_empty",
		},
		{
			"description exists but name is an empty string",
			test_utils.CreatePointerOfString(""),
			test_utils.CreatePointerOfString("description"),
			123,
			"name_or_description_empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mItemRepository.
				EXPECT().
				Create(gomock.Any()).
				Times(0)

			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/projects/%d", tc.projectId), test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 400 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 400, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateItemResponseData]{
				Success: false,
				Error:   tc.error,
			}

			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestCreateItemHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := CreateItemHandler(mItemRepository)

	router := httprouter.New()
	router.POST("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
		projectId   int
	}{
		{
			"object repository returns error",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString("description"),
			123,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mItemRepository.
				EXPECT().
				Create(gomock.Eq(&models.Item{
					Name:        *tc.name,
					Description: *tc.description,
					ProjectID:   uint(tc.projectId),
				})).
				Return(errors.New("error_string"))

			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/projects/%d", tc.projectId), test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 500 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 500, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateItemResponseData]{
				Success: false,
				Error:   "error_string",
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestCreateItemHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mItemRepository := mock_repository.NewMockItemRepository(ctrl)

	handleFunc := CreateItemHandler(mItemRepository)

	router := httprouter.New()
	router.POST("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
		projectId   int
	}{
		{
			"success case",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString("description"),
			123,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mItemRepository.
				EXPECT().
				Create(gomock.Eq(&models.Item{
					Name:        *tc.name,
					Description: *tc.description,
					ProjectID:   uint(tc.projectId),
				})).
				Do(func(m *models.Item) {
					m.ID = 345
				})

			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/projects/%d", tc.projectId), test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateItemResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateItemResponseData]{
				Success: true,
				Data:    CreateItemResponseData{ItemID: 345},
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
