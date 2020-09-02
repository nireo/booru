package models

import "github.com/jinzhu/gorm"

type Board struct {
	gorm.Model
	Title string `json:"title"`
	Link  string `json:"link"`
	UUID  string `json:"uuid"`
	Posts []Post
}
