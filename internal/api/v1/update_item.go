package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UpdateItemRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateItemResponseData struct{}

func UpdateItemHandler(repository repository.ItemRepository) httprouter.Handle {
	apiName := "update_item"
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody UpdateItemRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := utils.GenericResponse[UpdateItemResponseData]{Success: false, Error: "invalid_request_format"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		if requestBody.Name == nil && requestBody.Description == nil {
			responseBody := utils.GenericResponse[UpdateItemResponseData]{Success: false, Error: "name_and_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		fields := make(map[string]interface{})
		if requestBody.Name != nil && *requestBody.Name != "" {
			fields["name"] = requestBody.Name
		}

		if requestBody.Description != nil && *requestBody.Description != "" {
			fields["description"] = requestBody.Description
		}

		if len(fields) == 0 {
			responseBody := utils.GenericResponse[UpdateItemResponseData]{Success: false, Error: "name_and_description_empty"}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		err = repository.Update(projectId, itemId, fields)
		if err != nil {
			responseBody := utils.GenericResponse[UpdateItemResponseData]{Success: false, Error: "internal_db_error"}

			httpCode := 500
			if errors.Is(err, gorm.ErrRecordNotFound) {
				responseBody.Error = "item_or_project_not_found"
				httpCode = 404
			}

			utils.ReturnErrorResponse(r.Context(), w, httpCode, responseBody, err)
			utils.IncrementApiErrorMetric(apiName, responseBody.Error)
			return
		}

		responseBody := utils.GenericResponse[UpdateItemResponseData]{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
		utils.IncrementApiSuccessMetric(apiName)
	}
}
