package rest_api

import (
	"encoding/json"
	"net/http"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

func ServeBoardJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()
	var boards []models.Board
	if err := db.Find(&boards).Error; err != nil {
		http.Error(w, "Boards not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(boards)
}
