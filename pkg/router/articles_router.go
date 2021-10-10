package router

import (
	"github.com/6156-DonaldDuck/articles/pkg/config"
	"github.com/6156-DonaldDuck/articles/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func InitRouter() {
	r := gin.Default()
	r.GET("/articles", ListAllArticles)
	r.GET("/articles/:articleId", GetArticleByArticleId)
	r.Run(":" + config.Configuration.Port)
}

func ListAllArticles(c *gin.Context) {
	articles, err := service.ListAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, articles)
	}
}

func GetArticleByArticleId(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		log.Errorf("[router.GetArticleByArticleId] failed to parse article id %v, err=%v\n", articleIdStr, err)
		c.JSON(http.StatusBadRequest, "invalid article id")
		return
	}
	article, err := service.GetArticleByArticleId(uint(articleId))
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, article)
	}
}