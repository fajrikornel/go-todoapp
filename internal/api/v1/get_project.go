package v1

import (
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type GetProjectResponseBody struct {
	Success bool            `json:"success"`
	Project ProjectResponse `json:"project,omitempty"`
}

type ProjectResponse struct {
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
			responseBody := GetProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		logging.Infof(r.Context(), "Success getting project: %v", project)
		responseBody := buildGetProjectResponseBody(project)
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}

func buildGetProjectResponseBody(project *models.Project) GetProjectResponseBody {
	itemResponses := make([]ItemResponse, len(project.Items))
	for i, v := range project.Items {
		itemResponses[i] = ItemResponse{
			ItemID: v.ID,
			Name:   v.Name,
		}
	}

	responseBody := GetProjectResponseBody{Success: true, Project: ProjectResponse{
		ProjectID:   project.ID,
		Name:        project.Name,
		Description: project.Description,
		Items:       itemResponses,
	}}
	return responseBody
}
