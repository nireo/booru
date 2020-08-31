package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	AdminAccess      bool       `json:"adminAccess"`
	Port             string     `json:"port"`
	DatabaseName     string     `json:"databaseName"`
	DatabaseHost     string     `json:"databaseHost"`
	DatabasePort     string     `json:"databasePort"`
	DatabaseUser     string     `json:"databaseUser"`
	AdminCredentials AdminLogin `json:"adminLogin"`
}

type AdminLogin struct {
	Username string
	Password string
}

var conf Configuration

func LoadConfiguration() {
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
}

func GetConfiguration() Configuration {
	return conf
}
