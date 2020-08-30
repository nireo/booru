package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	UUID          string
	CreatedBy     string
	PostID        uint
	Content       string
	FileExtension string
}
