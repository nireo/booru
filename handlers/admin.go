package handlers

import (
	"net/http"

	"github.com/nireo/booru/models"
	"github.com/nireo/upfi/lib"
)

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	switch r.Method {
	case http.MethodPost:
		newBoard := &models.Board{
			UUID:  "unique",
			Link:  r.FormValue("link"),
			Title: r.FormValue("title"),
		}

		db.Save(&newBoard)
		http.Redirect(w, r, "http://localhost:8080/", http.StatusMovedPermanently)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	switch r.Method {
	case http.MethodPost:
		var board models.Board
		if err := db.Where("uuid = ?", "unique").First(&board).Error; err != nil {
			http.Error(w, "Board not found", http.StatusNotFound)
			return
		}

		db.Delete(&board)
		http.Redirect(w, r, "http://localhost:8080/", http.StatusMovedPermanently)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
