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

func TestCreateProjectHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := CreateProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.POST("/v1/projects", handleFunc)

	testCases := []struct {
		name        *string
		description *string
		error       string
	}{
		{
			nil,
			nil,
			"name_or_description_empty",
		},
		{
			createPointerOfString("name"),
			nil,
			"name_or_description_empty",
		},
		{
			nil,
			createPointerOfString("description"),
			"name_or_description_empty",
		},
		{
			createPointerOfString(""),
			createPointerOfString(""),
			"name_or_description_empty",
		},
		{
			createPointerOfString("name"),
			createPointerOfString(""),
			"name_or_description_empty",
		},
		{
			createPointerOfString(""),
			createPointerOfString("description"),
			"name_or_description_empty",
		},
	}
	for _, tc := range testCases {
		t.Run(formatTestCaseDescription(tc.name, tc.description), func(t *testing.T) {
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
				t.Errorf("Unexpected HTTP return code. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
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
		name        *string
		description *string
	}{
		{
			createPointerOfString("name"),
			createPointerOfString("description"),
		},
	}
	for _, tc := range testCases {
		t.Run(formatTestCaseDescription(tc.name, tc.description), func(t *testing.T) {
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
				t.Errorf("Unexpected HTTP return code. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
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
		name        *string
		description *string
	}{
		{
			createPointerOfString("name"),
			createPointerOfString("description"),
		},
	}
	for _, tc := range testCases {
		t.Run(formatTestCaseDescription(tc.name, tc.description), func(t *testing.T) {
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
				t.Errorf("Unexpected HTTP return code. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func createPointerOfString(s string) *string {
	sPointer := &s
	return sPointer
}

func formatTestCaseDescription(name, description *string) string {
	nameString := "nil"
	if name != nil {
		nameString = *name
	}
	descriptionString := "nil"
	if description != nil {
		descriptionString = *description
	}

	return fmt.Sprintf("name:%s and description:%s", nameString, descriptionString)
}
