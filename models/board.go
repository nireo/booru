package models

import "github.com/jinzhu/gorm"

type Board struct {
	gorm.Model
	Title string
	Link  string
	UUID  string
	Posts []Post
}
