package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type CreateItemRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateItemResponseBody struct {
	Success bool `json:"success"`
	ItemID  uint `json:"item_id"`
}

func CreateItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody CreateItemRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := CreateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, err)
			return
		}

		if requestBody.Name == "" || requestBody.Description == "" {
			responseBody := CreateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 400, responseBody, errors.New("invalid_format"))
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		item := models.Item{
			Name:        requestBody.Name,
			Description: requestBody.Description,
			ProjectID:   uint(projectId),
		}

		err = repository.Create(&item)
		if err != nil {
			responseBody := CreateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(w, 500, responseBody, err)
			return
		}

		log.Printf("Success creating item: %v", item)
		responseBody := CreateItemResponseBody{Success: true, ItemID: item.ID}
		utils.ReturnSuccessResponse(w, responseBody)
	}
}
