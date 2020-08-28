package main

import (
	"net/http"

	"github.com/nireo/booru/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ServeHomepage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
