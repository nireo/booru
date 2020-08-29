package models

import "github.com/jinzhu/gorm"

// named commentpost since comment is already in golang :(
type Comment struct {
	gorm.Model
	UUID          string
	CreatedBy     string
	PostID        uint
	Content       string
	FileExtension string
}
