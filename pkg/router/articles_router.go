package router

import (
	"errors"
	docs "github.com/6156-DonaldDuck/articles/docs"
	"github.com/6156-DonaldDuck/articles/pkg/config"
	"github.com/6156-DonaldDuck/articles/pkg/model"
	"github.com/6156-DonaldDuck/articles/pkg/router/middleware"
	"github.com/6156-DonaldDuck/articles/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) // use customized cors middleware
	r.Use(middleware.Security())
	r.Use(middleware.Notification())

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
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageStr := c.DefaultQuery("page", "1")
	authorIdStr := c.DefaultQuery("author_id", "0")
	sectionIdStr := c.DefaultQuery("section_id", "0")

	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		log.Errorf("[router.ListAllArticles] failed to parse authorId %v, err=%v\n", authorIdStr, err)
		c.JSON(http.StatusBadRequest, "invalid authorId")
		return
	}
	sectionId, err := strconv.Atoi(sectionIdStr)
	if err != nil {
		log.Errorf("[router.ListAllArticles] failed to parse sectionId %v, err=%v\n", sectionIdStr, err)
		c.JSON(http.StatusBadRequest, "invalid sectionId")
		return
	}

	pageSize, errPageSize := strconv.Atoi(pageSizeStr)
	page, errPage := strconv.Atoi(pageStr)
	if errPageSize != nil {
		log.Errorf("[router.ListAllArticles] failed to parse page size %v, err=%v\n", pageSizeStr, errPageSize)
		c.JSON(http.StatusBadRequest, "invalid page size")
		return
	}
	if errPage != nil {
		log.Errorf("[router.ListAllArticles] failed to parse page %v, err=%v\n", pageStr, errPage)
		c.JSON(http.StatusBadRequest, "invalid page")
		return
	}

	articles, total, err := service.ListAllArticles((page - 1) * pageSize, pageSize, uint(authorId), uint(sectionId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, model.ListArticlesResponse{
			Articles: articles,
			Total: total,
			Page: page,
			PageSize: pageSize,
		})
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
// @Router /articles/{ID} [get]
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Error(err)
		}
	} else{
		c.JSON(http.StatusOK, article)
	}
}

// @Summary Create Article
// @Schemes
// @Description Create Article
// @Tags Articles
// @Accept json
// @Produce json
// @Param title formData string false "Title"
// @Param content formData string false "Content"
// @Param kind formData string false "Kind"
// @Success 200 {json} article ids
// @Failure 400 invalid article id
// @Router /articles/ [post]
func CreateArticle(c *gin.Context) {
	article := model.Article{}
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if article.ID != 0 {
		_, err := service.GetArticleByArticleId(article.ID)
		if err == nil {
			c.JSON(http.StatusUnprocessableEntity, "Duplicate key")
		}
	}
	articleId, err := service.CreateArticle(article)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusCreated, articleId)
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
// @Router /articles/{ID} [put]
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
// @Router /articles/{ID} [delete]
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
		c.JSON(http.StatusNoContent, "Successfully delete article with id "+ idStr)
	}
}
