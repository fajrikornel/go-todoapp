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

type DeleteItemResponseData struct{}

func DeleteItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		err := repository.Delete(projectId, itemId)
		if err != nil {
			responseBody := utils.GenericResponse[DeleteItemResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "item_or_project_not_found"
				httpCode = 400
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			return
		}

		responseBody := utils.GenericResponse[DeleteItemResponseData]{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
