package handlers

import (
	"net/http"
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	switch r.Method {
	case http.MethodPost:
		newBoard := &models.Board{
			UUID:  lib.UUID(),
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
	query := r.URL.Query()
	boardUUID := query.Get("board")

	switch r.Method {
	case http.MethodPost:
		var board models.Board
		if err := db.Where("uuid = ?", boardUUID).First(&board).Error; err != nil {
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

func ServeManagerPage(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	var boards models.Board
	if err := db.Find(&boards).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/admin_manage.html"))
	if err := tmpl.Execute(w, boards); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
