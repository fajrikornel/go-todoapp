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

func TestDeleteProjectHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := DeleteProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.DELETE("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName             string
		projectId            int
		returnedError        error
		expectedErrorMessage string
		expectedHttpCode     int
	}{
		{
			"project does not exist in database",
			123,
			gorm.ErrRecordNotFound,
			"project_not_found",
			404,
		},
		{
			"error while calling database",
			123,
			gorm.ErrUnsupportedDriver,
			"internal_db_error",
			500,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mProjectRepository.
				EXPECT().
				Delete(gomock.Eq(tc.projectId)).
				Return(tc.returnedError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/projects/%d", tc.projectId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.expectedHttpCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.expectedHttpCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[DeleteProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[DeleteProjectResponseData]{
				Success: false,
				Error:   tc.expectedErrorMessage,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}

func TestDeleteProjectHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := DeleteProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.DELETE("/v1/projects/:projectId", handleFunc)

	testCases := []struct {
		testName  string
		projectId int
	}{
		{
			"success deleting project",
			123,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mProjectRepository.
				EXPECT().
				Delete(gomock.Eq(tc.projectId)).
				Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/projects/%d", tc.projectId), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if 200 != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", 200, rr.Code)
			}

			var actualResponse utils.GenericResponse[DeleteProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)

			expectedResponse := utils.GenericResponse[DeleteProjectResponseData]{
				Success: true,
			}
			if !reflect.DeepEqual(expectedResponse, actualResponse) {
				t.Errorf("Unexpected HTTP response. Expected: %+v, actual: %+v", expectedResponse, actualResponse)
			}
		})
	}
}
