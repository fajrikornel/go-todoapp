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

type GetItemResponseBody struct {
	Success bool         `json:"success"`
	Item    *models.Item `json:"item"`
}

func GetItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		item, err := repository.FindByProjectIdAndItemId(projectId, itemId)
		if err != nil {
			responseBody := GetItemResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, err)
			return
		}

		log.Printf("Success getting item: %v", item)
		responseBody := GetItemResponseBody{Success: true, Item: item}
		utils.ReturnSuccessResponse(w, responseBody)
	}
}
