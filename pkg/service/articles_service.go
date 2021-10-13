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

func CreateArticle(article model.Article) (articleId uint) {
	result := db.DbConn.Create(&article)
	if result.Error != nil {
		log.Errorf("[service.CreateArticle] error occurred while creating article, err=%v\n", result.Error)
	} else {
		log.Infof("[service.CreateArticle] successfully created article with id %v, rows affected = %v\n", article.ID, result.RowsAffected)
	}
	return article.ID, result.Error
}

func UpdateArticle(updateInfo model.Article) (error){
	result := db.DbConn.Model(&updateInfo).Updates(updateInfo)
	if result.Error != nil {
		log.Errorf("[service.UpdateArticle] error occurred while updating article, err=%v\n", result.Error)
	} else {
		log.Infof("[service.UpdateArticle] successfully updated article with id %v, rows affected = %v\n", updateInfo.ID, result.RowsAffected)
	}
	return result.Error
}

func DeleteArticleById(articleId uint) (error) {
	article := model.Article{}
	result := db.DbConn.Delete(&article, articleId)
	if result.Error != nil {
		log.Errorf("[service.DeleteArticleById] error occurred while deleting article with id %v, err=%v\n", articleId, result.Error)
	} else {
		log.Infof("[service.DeleteArticleById] successfully deleted article with id %v, rows affected = %v\n", articleId, result.RowsAffected)
	}
	return result.Error
}