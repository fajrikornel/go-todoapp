package v1

import (
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type DeleteItemResponseBody struct {
	Success bool `json:"success"`
}

func DeleteItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		_, err := repository.FindByProjectIdAndItemId(projectId, itemId)
		if err != nil {
			responseBody := DeleteItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		err = repository.Delete(projectId, itemId)
		if err != nil {
			responseBody := DeleteItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		responseBody := DeleteItemResponseBody{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
