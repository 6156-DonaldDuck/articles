package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	AuthorId uint `json:"author_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	SectionId uint `json:"section_id"`
}

type DArticle struct {
	AuthorId uint
	Title string 
	Content string
	SectionId uint
}