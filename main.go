package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/nireo/booru/handlers"
	"github.com/nireo/booru/lib"
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
	// load configuration file
	var conf Configuration
	configurationFile, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer configurationFile.Close()

	inBytes, err := ioutil.ReadAll(configurationFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(inBytes, &conf)

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

	// setup http server and all the handlers
	http.HandleFunc("/board", handlers.GetPostsInBoard)
	http.HandleFunc("/post", handlers.CreateNewPost)
	http.HandleFunc("/", handlers.ServeHomepage)

	if conf.AdminAccess {
		http.HandleFunc("/board/create", handlers.CreateBoard)
		http.HandleFunc("/board/delete", handlers.DeleteBoard)
		http.HandleFunc("/manage", handlers.ServeManagerPage)
	}

	if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
		panic(err)
	}
}
