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
	AuthorId uint `json:"author_id" dynamodbav:"AuthorId"`
	Title string `json:"title" dynamodbav:"Title"`
	Content string `json:"content" dynamodbav:"Content"`
	SectionId uint `json:"section_id" dynamodbav:"SectionId"`
}