package handlers

import (
	"net/http"
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()
	store := lib.GetStore()
	session, _ := store.Get(r, "admin-auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
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
	store := lib.GetStore()
	session, _ := store.Get(r, "admin-auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	keys, ok := r.URL.Query()["post"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide post id", http.StatusBadRequest)
		return
	}
	db := lib.GetDatabase()
	store := lib.GetStore()
	session, _ := store.Get(r, "admin-auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var post models.Post
	if err := db.Where(&models.Post{UUID: keys[0]}).First(&post).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	db.Delete(&post)
}

func ServeManagerPage(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDatabase()

	var boards models.Board
	if err := db.Find(&boards).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	store := lib.GetStore()
	session, _ := store.Get(r, "admin-auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/admin_manage.html"))
	if err := tmpl.Execute(w, boards); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	store := lib.GetStore()
	session, _ := store.Get(r, "admin-auth")
	conf := lib.GetConfiguration()

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != conf.AdminCredentials.Username {
		http.Error(w, "Incorrect credentials", http.StatusForbidden)
		return
	}

	if password != conf.AdminCredentials.Password {
		http.Error(w, "Incorrect credentials", http.StatusForbidden)
		return
	}

	session.Values["authorized"] = true
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/", http.StatusMovedPermanently)
}
