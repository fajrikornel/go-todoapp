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

type GetItemResponseData struct {
	ItemID      uint   `json:"item_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func GetItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		item, err := repository.FindByProjectIdAndItemId(projectId, itemId)
		if err != nil {
			responseBody := utils.GenericResponse[GetItemResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "item_or_project_not_found"
				httpCode = 404
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			return
		}

		responseBody := utils.GenericResponse[GetItemResponseData]{Success: true, Data: GetItemResponseData{
			ItemID:      item.ID,
			Name:        item.Name,
			Description: item.Description,
		}}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
