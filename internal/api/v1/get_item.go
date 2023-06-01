package v1

import (
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type GetItemResponseBody struct {
	Success bool                 `json:"success"`
	Item    DetailedItemResponse `json:"item"`
}

type DetailedItemResponse struct {
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
			responseBody := GetItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		logging.Infof(r.Context(), "Success getting item: %v", item)
		responseBody := GetItemResponseBody{Success: true, Item: DetailedItemResponse{
			ItemID:      item.ID,
			Name:        item.Name,
			Description: item.Description,
		}}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
