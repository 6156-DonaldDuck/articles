package router

import (
	"net/http"
	"strconv"

	docs "github.com/6156-DonaldDuck/articles/docs"
	"github.com/6156-DonaldDuck/articles/pkg/config"
	"github.com/6156-DonaldDuck/articles/pkg/model"
	"github.com/6156-DonaldDuck/articles/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = config.Configuration.BaseURL
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group(config.Configuration.BaseURL)
	{
		apiv1.GET("/articles", ListAllArticles)
		apiv1.GET("/articles/:articleId", GetArticleByArticleId)
		apiv1.POST("/articles", CreateArticle)
		apiv1.PUT("/articles/:articleId", UpdateArticleById)
		apiv1.DELETE("/articles/:articleId", DeleteArticleById)
	}

	r.Run(":" + config.Configuration.Port)
}

// @BasePath /api/v1

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
// @Param ID path int true "the id of a specfic article"
// @Success 200 {json} article
// @Failure 400 invalid article id
// @Router /articles/{articleId} [get]
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

// @Summary Create Article
// @Schemes
// @Description Create Article
// @Tags Articles
// @Accept json
// @Produce json
// @Param ID formData int true "the id of a specfic article"
// @Param title formData string false "Title"
// @Param content formData string false "Content"
// @Param kind formData string false "Kind"
// @Success 200 {json} article id
// @Failure 400 invalid article id
// @Router /articles/{articleId} [post]
func CreateArticle(c *gin.Context) {
	article := model.Article{}
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	articleId, err := service.CreateArticle(article)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, articleId)
	}
}

// @Summary Update Article By Article Id
// @Schemes
// @Description Update article by article id
// @Tags Articles
// @Accept json
// @Produce json
// @Param ID path int true "the id of a specfic article"
// @Success 200 {json} update successfully
// @Failure 400 invalid article id
// @Router /articles/ [put]
func UpdateArticleById(c *gin.Context) {
	idStr := c.Param("articleId")
	articleId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.UpdateArticleById] failed to parse article id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid article id")
		return
	}
	updateInfo := model.Article{}
	if err := c.ShouldBind(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	updateInfo.ID = uint(articleId)
	err = service.UpdateArticle(updateInfo)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "update successfully")
	}
}

// @Summary Delete Article By Article Id
// @Schemes
// @Description Delete article by article id
// @Tags Articles
// @Accept json
// @Produce json
// @Param ID header int true "the id of a specfic article"
// @Success 200 {json} delete successfully
// @Failure 400 invalid article id
// @Router /articles/ [delete]
func DeleteArticleById(c *gin.Context) {
	idStr := c.Param("articleId")
	articleId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.DeleteArticleById] failed to parse article id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid article id")
		return
	}
	err = service.DeleteArticleById(uint(articleId))
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "Successfully delete article with id "+idStr)
	}
}
