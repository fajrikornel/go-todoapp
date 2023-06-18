package v1

import (
	"encoding/json"
	"errors"
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

func TestCreateProjectHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := CreateProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.POST("/v1/projects", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
		error       string
	}{
		{
			"name and description does not exist",
			nil,
			nil,
			"name_or_description_empty",
		},
		{
			"name exists but description does not exist",
			test_utils.CreatePointerOfString("name"),
			nil,
			"name_or_description_empty",
		},
		{
			"description exists but name does not exist",
			nil,
			test_utils.CreatePointerOfString("description"),
			"name_or_description_empty",
		},
		{
			"name and description are empty strings",
			test_utils.CreatePointerOfString(""),
			test_utils.CreatePointerOfString(""),
			"name_or_description_empty",
		},
		{
			"name exists but description is an empty string",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString(""),
			"name_or_description_empty",
		},
		{
			"description exists but name is an empty string",
			test_utils.CreatePointerOfString(""),
			test_utils.CreatePointerOfString("description"),
			"name_or_description_empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mProjectRepository.
				EXPECT().
				Create(gomock.Any()).
				Times(0)

			req := httptest.NewRequest("POST", "/v1/projects", test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 400 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 400, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   tc.error,
			}

			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestCreateProjectHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := CreateProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.POST("/v1/projects", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
	}{
		{
			"object repository returns error",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString("description"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mProjectRepository.
				EXPECT().
				Create(gomock.Eq(&models.Project{
					Name:        *tc.name,
					Description: *tc.description,
				})).
				Return(errors.New("error_string"))

			req := httptest.NewRequest("POST", "/v1/projects", test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 500 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 500, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "error_string",
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestCreateProjectHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := CreateProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.POST("/v1/projects", handleFunc)

	testCases := []struct {
		testName    string
		name        *string
		description *string
	}{
		{
			"success case",
			test_utils.CreatePointerOfString("name"),
			test_utils.CreatePointerOfString("description"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			mProjectRepository.
				EXPECT().
				Create(gomock.Eq(&models.Project{
					Name:        *tc.name,
					Description: *tc.description,
				})).
				Do(func(m *models.Project) {
					m.ID = 123
				})

			req := httptest.NewRequest("POST", "/v1/projects", test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[CreateProjectResponseData]{
				Success: true,
				Data:    CreateProjectResponseData{ProjectID: 123},
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
