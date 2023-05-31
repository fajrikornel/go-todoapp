package v1

import (
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type GetProjectResponseBody struct {
	Success bool            `json:"success"`
	Project *models.Project `json:"project"`
}

func GetProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		project, err := repository.FindById(projectId)
		if err != nil {
			responseBody := GetProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, err)
			return
		}

		log.Printf("Success getting project: %v", project)
		responseBody := GetProjectResponseBody{Success: true, Project: project}
		utils.ReturnSuccessResponse(w, responseBody)
	}
}
