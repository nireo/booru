package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

type BoardPage struct {
	Board models.Board
	Posts []models.Post
}

func GetPostsInBoard(w http.ResponseWriter, r *http.Request) {
	// fetch board link
	keys, ok := r.URL.Query()["board"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	// load board and posts
	db := lib.GetDatabase()
	var board models.Board
	if err := db.Where(&models.Board{Link: keys[0]}).First(&board).Error; err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	// don't return error, since board can be empty
	var posts []models.Post
	db.Model(&board).Related(&posts)

	fmt.Println(posts)

	boardPage := &BoardPage{
		Board: board,
		Posts: posts,
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/posts.html"))
	if err := tmpl.Execute(w, boardPage); err != nil {
		http.Error(w, "Internal server error", http.StatusNotFound)
		return
	}
}

func CreateNewPost(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	// fetch board link
	keys, ok := r.URL.Query()["board"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	// make sure the file size is under 10 mbs
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if the board exists
	var board models.Board
	if err := db.Where(&models.Board{Link: keys[0]}).First(&board).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPost := &models.Post{
		UUID:          lib.UUID(),
		CreatedBy:     "Anonymous",
		Content:       r.FormValue("content"),
		BoardID:       board.ID,
		FileExtension: filepath.Ext(handler.Filename),
	}

	// save the file
	defer file.Close()
	filePath := fmt.Sprintf("./images/%s%s", newPost.UUID, newPost.FileExtension)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// finally create the database entry, and redirect the user to the board page
	db.NewRecord(newPost)
	db.Create(newPost)
	http.Redirect(
		w, r,
		fmt.Sprintf("http://localhost:8080/board/?board=%s", board.Link),
		http.StatusMovedPermanently)
}
