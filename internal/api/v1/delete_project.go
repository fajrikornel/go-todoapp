package v1

import (
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type DeleteProjectResponseBody struct {
	Success bool `json:"success"`
}

func DeleteProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))

		_, err := repository.FindById(projectId)
		if err != nil {
			responseBody := DeleteProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		err = repository.Delete(projectId)
		if err != nil {
			responseBody := DeleteProjectResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		responseBody := DeleteProjectResponseBody{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}