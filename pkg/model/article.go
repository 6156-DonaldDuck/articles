package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	AuthorId uint `json:"author_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Kind string `json:"kind"`
}