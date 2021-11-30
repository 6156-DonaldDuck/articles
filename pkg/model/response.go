package model

type ListArticlesResponse struct {
	Articles []Article `json:"articles"`
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListArticlesResponseD struct {
	Articles []DArticle `json:"articles"`
}