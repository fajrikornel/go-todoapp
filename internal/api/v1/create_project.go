package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CreateProjectRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateProjectResponseData struct {
	ProjectID uint `json:"project_id,omitempty"`
}

func CreateProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody CreateProjectRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := utils.GenericResponse[CreateProjectResponseData]{Success: false, Error: err.Error()}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		if requestBody.Name == nil || requestBody.Description == nil {
			responseBody := utils.GenericResponse[CreateProjectResponseData]{Success: false, Error: "name_or_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_or_description_empty"))
			return
		}

		if *requestBody.Name == "" || *requestBody.Description == "" {
			responseBody := utils.GenericResponse[CreateProjectResponseData]{Success: false, Error: "name_or_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_or_description_empty"))
			return
		}

		project := models.Project{
			Name:        *requestBody.Name,
			Description: *requestBody.Description,
		}

		err = repository.Create(&project)
		if err != nil {
			responseBody := utils.GenericResponse[CreateProjectResponseData]{Success: false, Error: err.Error()}
			utils.ReturnErrorResponse(r.Context(), w, 500, responseBody, err)
			return
		}

		responseBody := utils.GenericResponse[CreateProjectResponseData]{Success: true, Data: CreateProjectResponseData{ProjectID: project.ID}}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
