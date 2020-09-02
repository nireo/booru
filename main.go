package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/nireo/booru/handlers"
	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/middleware"
	"github.com/nireo/booru/models"
)

type Configuration struct {
	AdminAccess  bool   `json:"adminAccess"`
	Port         string `json:"port"`
	DatabaseName string `json:"databaseName"`
	DatabaseHost string `json:"databaseHost"`
	DatabasePort string `json:"databasePort"`
	DatabaseUser string `json:"databaseUser"`
}

func main() {
	lib.LoadConfiguration()
	conf := lib.GetConfiguration()

	// database connection and model migration
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		conf.DatabaseHost, conf.DatabasePort, conf.DatabaseUser, conf.DatabaseName)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	models.MigrateModels(db)
	defer db.Close()
	lib.SetDatabase(db)

	loggingMiddleware := middleware.NestedMiddleware(
		middleware.LogRequest, middleware.TraceRequest)

	// setup http server and all the handlers
	// (rest api mode only hosts routes that give out json output, and doesn't serve routes that use templates)
	// images are still served in every mode
	if conf.RestAPIMode {

	} else {
		http.HandleFunc("/board", loggingMiddleware(handlers.GetPostsInBoard))
		http.HandleFunc("/post", loggingMiddleware(handlers.GetSinglePost))
		http.HandleFunc("/post/create", loggingMiddleware(handlers.CreateNewPost))
		http.HandleFunc("/reply", loggingMiddleware(handlers.CreateComment))
		http.HandleFunc("/", handlers.ServeHomepage)
	}
	fs := http.FileServer(http.Dir("images/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	if conf.AdminAccess {
		http.HandleFunc("/board/create", handlers.CreateBoard)
		http.HandleFunc("/board/delete", handlers.DeleteBoard)
		http.HandleFunc("/manage", handlers.ServeManagerPage)
	}

	if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
		panic(err)
	}
}
