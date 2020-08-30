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

type SinglePostPage struct {
	Post     models.Post
	Comments []models.Comment
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	// fetch board link
	keys, ok := r.URL.Query()["post"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	// make sure the file size is under 10 mbs
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var post models.Post
	if err := db.Where(&models.Post{UUID: keys[0]}).First(&post).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// check if the comment has a file, since comments don't have to have images
	withFile := true
	file, handler, err := r.FormFile("file")
	if err != nil {
		withFile = false
	}

	newComment := &models.Comment{
		UUID:      lib.UUID(),
		CreatedBy: "Anonymous",
		Content:   r.FormValue("content"),
		PostID:    post.ID,
	}

	if withFile {
		newComment.FileExtension = filepath.Ext(handler.Filename)
		// save the file
		defer file.Close()
		filePath := fmt.Sprintf("./images/%s%s", newComment.UUID, newComment.FileExtension)
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
	}

	// finally create the database entry, and redirect the user to the board page
	db.NewRecord(newComment)
	db.Create(newComment)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully answered"))
}

func GetSinglePost(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	keys, ok := r.URL.Query()["post"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide post name", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := db.Where(&models.Post{UUID: keys[0]}).First(&post).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, _ := post.GetComments()

	singlePostData := &SinglePostPage{
		Comments: comments,
		Post:     post,
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/single_post.html"))
	if err := tmpl.Execute(w, singlePostData); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
