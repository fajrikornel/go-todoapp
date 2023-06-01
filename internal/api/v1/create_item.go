package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type CreateItemRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateItemResponseData struct {
	ItemID uint `json:"item_id,omitempty"`
}

func CreateItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody CreateItemRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := utils.GenericResponse[CreateItemResponseData]{Success: false, Error: err.Error()}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		if requestBody.Name == nil || requestBody.Description == nil {
			responseBody := utils.GenericResponse[CreateItemResponseData]{Success: false, Error: "name_or_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_or_description_empty"))
			return
		}

		if *requestBody.Name == "" || *requestBody.Description == "" {
			responseBody := utils.GenericResponse[CreateItemResponseData]{Success: false, Error: "name_or_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_or_description_empty"))
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		item := models.Item{
			Name:        *requestBody.Name,
			Description: *requestBody.Description,
			ProjectID:   uint(projectId),
		}

		err = repository.Create(&item)
		if err != nil {
			responseBody := utils.GenericResponse[CreateItemResponseData]{Success: false, Error: err.Error()}
			utils.ReturnErrorResponse(r.Context(), w, 500, responseBody, err)
			return
		}

		responseBody := utils.GenericResponse[CreateItemResponseData]{Success: true, Data: CreateItemResponseData{ItemID: item.ID}}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
