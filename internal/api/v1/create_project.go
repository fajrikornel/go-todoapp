package v1

import (
	"encoding/json"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/models"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type CreateProjectRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectResponseBody struct {
	Success bool `json:"success"`
}

func CreateProjectHandler(store *db.SqlStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var requestBody CreateProjectRequestBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			log.Printf("Bad request: %v", err.Error())
			responseBody := CreateProjectResponseBody{Success: false}
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(responseBody)
			return
		}

		if requestBody.Name == "" || requestBody.Description == "" {
			log.Printf("Bad request: empty name or description")
			responseBody := CreateProjectResponseBody{Success: false}
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(responseBody)
			return
		}

		project := models.Project{
			Name:        requestBody.Name,
			Description: requestBody.Description,
		}

		err = store.Create(&project)
		if err != nil {
			log.Printf("Error saving to DB: %v", err.Error())
			responseBody := CreateProjectResponseBody{Success: false}
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(responseBody)
			return
		}

		log.Printf("Success creating project: %v", project)
		responseBody := CreateProjectResponseBody{Success: true}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
	}
}
