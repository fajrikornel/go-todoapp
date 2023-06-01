package v1

import (
	"encoding/json"
	"errors"
	"github.com/fajrikornel/go-todoapp/internal/api/utils"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type UpdateItemRequestBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateItemResponseBody struct {
	Success bool `json:"success"`
}

func UpdateItemHandler(repository repository.ItemRepository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody UpdateItemRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseBody := UpdateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, err)
			return
		}

		projectId, _ := strconv.Atoi(p.ByName("projectId"))
		itemId, _ := strconv.Atoi(p.ByName("itemId"))

		_, err = repository.FindByProjectIdAndItemId(projectId, itemId)
		if err != nil {
			responseBody := UpdateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("no_matching_project_and_item"))
			return
		}

		if requestBody.Name == nil && requestBody.Description == nil {
			responseBody := UpdateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 400, responseBody, errors.New("name_and_description_empty"))
			return
		}

		fields := make(map[string]interface{})
		if requestBody.Name != nil {
			fields["name"] = requestBody.Name
		}

		if requestBody.Description != nil {
			fields["description"] = requestBody.Description
		}

		err = repository.Update(projectId, itemId, fields)
		if err != nil {
			responseBody := UpdateItemResponseBody{Success: false}
			utils.ReturnErrorResponse(r.Context(), w, 500, responseBody, err)
			return
		}

		responseBody := UpdateItemResponseBody{Success: true}
		utils.ReturnSuccessResponse(r.Context(), w, responseBody)
	}
}
