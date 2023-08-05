package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UpdateProjectRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateProjectResponseData struct{}

func UpdateProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	apiName := "update_project"
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody UpdateProjectRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := utils.GenericResponse[UpdateProjectResponseData]{Success: false, Error: "invalid_request_format"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))

		if requestBody.Name == nil && requestBody.Description == nil {
			responseBody := utils.GenericResponse[UpdateProjectResponseData]{Success: false, Error: "name_and_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		fields := make(map[string]interface{})
		if requestBody.Name != nil && *requestBody.Name != "" {
			fields["name"] = requestBody.Name
		}

		if requestBody.Description != nil && *requestBody.Description != "" {
			fields["description"] = requestBody.Description
		}

		if len(fields) == 0 {
			responseBody := utils.GenericResponse[UpdateProjectResponseData]{Success: false, Error: "name_and_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		err = repository.Update(projectId, fields)
		if err != nil {
			responseBody := utils.GenericResponse[UpdateProjectResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "project_not_found"
				httpCode = 404
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		responseBody := utils.GenericResponse[UpdateProjectResponseData]{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
		utils.IncrementApiSuccessMetric(apiName)
	}
}
