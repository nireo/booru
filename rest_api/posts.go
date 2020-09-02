package rest_api

import (
	"encoding/json"
	"net/http"

	"github.com/nireo/booru/models"
	"github.com/nireo/upfi/lib"
)

func ServePostsInBoardJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// fetch board name
	keys, ok := r.URL.Query()["board"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()

	var board models.Board
	if err := db.Where(&models.Board{UUID: keys[0]}).First(&board).Error; err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	var posts []models.Post
	db.Model(&board).Related(&posts)

	for index := range posts {
		comments, _ := posts[index].GetComments()
		posts[index].Comments = comments
	}

	json.NewEncoder(w).Encode(posts)
}
