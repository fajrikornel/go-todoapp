package v1

import (
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type GetProjectResponseData struct {
	ProjectID   uint           `json:"project_id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Items       []ItemResponse `json:"items,omitempty"`
}

type ItemResponse struct {
	ItemID uint   `json:"item_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

func GetProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		project, err := repository.FindById(projectId)
		if err != nil {
			responseBody := utils.GenericResponse[GetProjectResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "project_not_found"
				httpCode = 400
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			return
		}

		responseBody := buildGetProjectResponseBody(project)
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}

func buildGetProjectResponseBody(project *models.Project) utils.GenericResponse[GetProjectResponseData] {
	itemResponses := make([]ItemResponse, len(project.Items))
	for i, v := range project.Items {
		itemResponses[i] = ItemResponse{
			ItemID: v.ID,
			Name:   v.Name,
		}
	}

	responseBody := utils.GenericResponse[GetProjectResponseData]{Success: true, Data: GetProjectResponseData{
		ProjectID:   project.ID,
		Name:        project.Name,
		Description: project.Description,
		Items:       itemResponses,
	}}
	return responseBody
}
