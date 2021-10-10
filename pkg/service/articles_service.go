package service

import (
	"github.com/6156-DonaldDuck/articles/pkg/db"
	"github.com/6156-DonaldDuck/articles/pkg/model"
	log "github.com/sirupsen/logrus"
)

func ListAllArticles() ([]model.Article, error) {
	var articles []model.Article
	result := db.DbConn.Find(&articles)
	if result.Error != nil {
		log.Errorf("[service.ListAllArticles] error occurred while listing articles, err=%v\n", result.Error)
	} else {
		log.Infof("[service.ListAllArticles] successfully listed articles, rows affected = %v\n", result.RowsAffected)
	}
	return articles, result.Error
}

func GetArticleByArticleId(articleId uint) (model.Article, error) {
	article := model.Article{}
	result := db.DbConn.First(&article, articleId)
	if result.Error != nil {
		log.Errorf("[service.GetArticleByArticleId] error occurred while getting article with id %v, err=%v\n", articleId, result.Error)
	} else {
		log.Infof("[service.GetArticleByArticleId] successfully got article with id %v, rows affected = %v\n", articleId, result.RowsAffected)
	}
	return article, result.Error
}