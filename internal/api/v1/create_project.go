package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type CreateProjectRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectResponseBody struct {
	Success   bool `json:"success"`
	ProjectID uint `json:"project_id"`
}

func CreateProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody CreateProjectRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := CreateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, err)
			return
		}

		if requestBody.Name == "" || requestBody.Description == "" {
			responseBody := CreateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, errors.New("invalid_format"))
			return
		}

		project := models.Project{
			Name:        requestBody.Name,
			Description: requestBody.Description,
		}

		err = repository.Create(&project)
		if err != nil {
			responseBody := CreateProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 500, responseBody, err)
			return
		}

		log.Printf("Success creating project: %v", project)
		responseBody := CreateProjectResponseBody{Success: true, ProjectID: project.ID}
		utils.ReturnSuccessResponse(w, responseBody)
	}
}
