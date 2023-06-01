package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type UpdateProjectRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateProjectResponseBody struct {
	Success bool `json:"success"`
}

func UpdateProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody UpdateProjectRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := UpdateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))

		if requestBody.Name == nil && requestBody.Description == nil {
			responseBody := UpdateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
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
			responseBody := UpdateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			return
		}

		err = repository.Update(projectId, fields)
		if err != nil {
			responseBody := UpdateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 500, responseBody, err)
			return
		}

		responseBody := UpdateProjectResponseBody{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
