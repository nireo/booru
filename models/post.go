package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nireo/booru/lib"
)

type Post struct {
	gorm.Model
	UUID          string `json:"uuid"`
	CreatedBy     string `json:"created_by"`
	BoardID       uint
	Comments      []Comment `json:"comments"`
	Content       string    `json:"content"`
	FileExtension string    `json:"file_extension"`
}

func GetPostComments(id string) ([]Comment, error) {
	db := lib.GetDatabase()

	var post Post
	var comments []Comment
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		return comments, err
	}

	if err := db.Model(&post).Related(&comments).Error; err != nil {
		return comments, err
	}

	return comments, nil
}

func (post *Post) GetComments() ([]Comment, error) {
	db := lib.GetDatabase()
	var comments []Comment
	if err := db.Model(&post).Related(&comments).Error; err != nil {
		return comments, err
	}

	return comments, nil
}
