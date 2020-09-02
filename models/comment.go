package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	UUID          string `json:"uuid"`
	CreatedBy     string `json:"created_by"`
	PostID        uint
	Content       string `json:"content"`
	FileExtension string `json:"file_extension"`
}
