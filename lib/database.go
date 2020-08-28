package lib

import "github.com/jinzhu/gorm"

var db *gorm.DB

func SetDatabase(database *gorm.DB) {
	db = database
}

func GetDatabase() *gorm.DB {
	return db
}
