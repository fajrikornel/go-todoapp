package v1

import (
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type DeleteProjectResponseData struct{}

func DeleteProjectHandler(repository repository.ProjectRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))

		err := repository.Delete(projectId)
		if err != nil {
			responseBody := utils.GenericResponse[DeleteProjectResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "project_not_found"
				httpCode = 404
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			return
		}

		responseBody := utils.GenericResponse[DeleteProjectResponseData]{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
