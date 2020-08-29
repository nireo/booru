package handlers

import (
	"net/http"
	"text/template"

	"github.com/nireo/booru/models"
	"github.com/nireo/upfi/lib"
)

func GetPostsInBoard(w http.ResponseWriter, r *http.Request) {
	// fetch board link
	query := r.URL.Query()
	boardName := query.Get("board")

	// load board and posts
	db := lib.GetDatabase()
	var board models.Board
	if err := db.Where("link = ?", boardName).First(&board); err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	// don't return error, since board can be empty
	var posts []models.Post
	db.Model(&board).Related(&posts)

	tmpl := template.Must(template.ParseFiles("./templates/pages/posts.html"))
	if err := tmpl.Execute(w, posts); err != nil {
		http.Error(w, "Internal server error", http.StatusNotFound)
		return
	}
}
