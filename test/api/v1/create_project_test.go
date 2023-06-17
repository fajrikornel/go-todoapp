package v1

import (
	"encoding/json"
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

func TestCreateProjectHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mProjectRepository := mock_repository.NewMockProjectRepository(ctrl)

	handleFunc := CreateProjectHandler(mProjectRepository)

	router := httprouter.New()
	router.POST("/v1/projects", handleFunc)

	testCases := []struct {
		name           *string
		description    *string
		shouldSaveToDb bool
		responseCode   int
		response       utils.GenericResponse[CreateProjectResponseData]
	}{
		{nil,
			nil,
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{createPointerOfString("name"),
			nil,
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{nil,
			createPointerOfString("description"),
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{createPointerOfString(""),
			createPointerOfString(""),
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{createPointerOfString("name"),
			createPointerOfString(""),
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{createPointerOfString(""),
			createPointerOfString("description"),
			false,
			400,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: false,
				Error:   "name_or_description_empty",
			},
		},
		{createPointerOfString("name"),
			createPointerOfString("description"),
			true,
			200,
			utils.GenericResponse[CreateProjectResponseData]{
				Success: true,
				Data:    CreateProjectResponseData{ProjectID: 123},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(formatTitle(tc.name, tc.description), func(t *testing.T) {
			requestBody := CreateProjectRequestBody{
				Name:        tc.name,
				Description: tc.description,
			}

			if tc.shouldSaveToDb {
				mProjectRepository.
					EXPECT().
					Create(gomock.Eq(&models.Project{
						Name:        *tc.name,
						Description: *tc.description,
					})).
					Do(func(m *models.Project) {
						m.ID = 123
					})
			}

			req := httptest.NewRequest("POST", "/v1/projects", test_utils.ConvertStructToIoReader(requestBody))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tc.responseCode != rr.Code {
				t.Errorf("Unexpected HTTP return code. Expected: %d, actual: %d", tc.responseCode, rr.Code)
			}

			var actualResponse utils.GenericResponse[CreateProjectResponseData]
			json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			if !reflect.DeepEqual(tc.response, actualResponse) {
				t.Errorf("Unexpected HTTP return code. Expected: %+v, actual: %+v", tc.response, actualResponse)
			}
		})
	}
}

func createPointerOfString(s string) *string {
	sPointer := &s
	return sPointer
}

func formatTitle(name, description *string) string {
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
