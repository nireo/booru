package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

type HomepageData struct {
	Boards []models.Board
}

func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()
	// get all the boards
	var boards []models.Board
	if err := db.Find(&boards).Error; err != nil {
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	fmt.Println(boards)

	homepageData := &HomepageData{
		Boards: boards,
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/boards.html"))
	if err := tmpl.Execute(w, homepageData); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
