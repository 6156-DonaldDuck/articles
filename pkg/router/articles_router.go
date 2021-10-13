package router

import (
	docs "github.com/6156-DonaldDuck/articles/docs"
	"github.com/6156-DonaldDuck/articles/pkg/config"
	"github.com/6156-DonaldDuck/articles/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
)

func InitRouter() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/articles", ListAllArticles)
	r.GET("/articles/:articleId", GetArticleByArticleId)
	r.Run(":" + config.Configuration.Port)
}

// @BasePath /


// @Summary List All Articles
// @Schemes
// @Description List all articles
// @Tags Articles
// @Accept json
// @Produce json
// @Success 200 {json} articles
// @Failure 500 internal server error
// @Router /articles [get]
func ListAllArticles(c *gin.Context) {
	articles, err := service.ListAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, articles)
	}
}


// @Summary Get Article By Article Id
// @Schemes
// @Description Get article by article id
// @Tags Articles
// @Accept json
// @Produce json
// @Param ID query int true "the id of a specfic article"
// @Success 200 {json} article
// @Failure 400 invalid article id
// @Router /articles/ [get]
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