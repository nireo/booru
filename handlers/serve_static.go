package handlers

import (
	"net/http"
	"text/template"
)

func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/pages/boards.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
