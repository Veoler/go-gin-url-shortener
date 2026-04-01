package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Veoler/go-gin-url-shortener/data"
	"github.com/Veoler/go-gin-url-shortener/models"
	"strconv"
)

func main() {
	r := gin.Default()

	r.GET("/links", GetAllLinks)
	r.GET("/links/:slug", GetLinkBySlug)
	r.POST("/links", AddLink)
	r.PATCH("/links/:id", UpdateLink)
	r.DELETE("/links/:id", DeleteLinkByID)
	r.DELETE("/links/slug/:slug", DeleteLinkBySlug)
	r.GET("/r/:slug", RedirectLink)

	r.Run(":8080")
}

func GetAllLinks(c *gin.Context) {
	linksList, err := data.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, linksList)
}

func GetLinkBySlug(c *gin.Context) {
	slug := c.Param("slug")
	linksSlug, err := data.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, linksSlug)
}

func UpdateLink (c *gin.Context) {
	id, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "ID должен быть числом"})
    	return
	}

	var input models.Link 
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при изменении"})
		return
	}

	hasURL := input.URL != nil && *input.URL != ""
	hasSlug := input.Slug != nil && *input.Slug != ""

	switch {
	// Проверка наличия хотя бы одного поля
	case !hasURL && !hasSlug:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поля url или slug обязательны(что-то одно)"})
		return

	// Проверка формата URL
	case data.IsInvalidURL(input.URL):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат URL"})
		return

	// Проверка дублирования slug 
	case hasSlug && data.IsSlugExistByOther(*input.Slug, id):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ссылка с таким slug уже существует"})
		return

	default:
		updated, err := data.Update(id, input)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	}
}

func AddLink (c *gin.Context) {
	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при добавлении: некорректные данные"})
		return
	}

	switch {
	// Проверка наличия url и slug
	case input.URL == nil || *input.URL == "" || input.Slug == nil || *input.Slug == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поля url и slug обязательны"})
		return
	
	// Проверка формата URL
	case data.IsInvalidURL(input.URL):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат URL"})
		return

	// Проверка дублирования slug 
	case data.IsSlugExist(*input.Slug):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ссылка с таким slug уже существует"})
		return

	default:
		c.JSON(http.StatusCreated, data.Add(input))
	}
}

func DeleteLinkByID(c *gin.Context) {
	id, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "ID должен быть числом"})
    	return
	}

	err := data.DeleteID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func DeleteLinkBySlug(c *gin.Context) {
	slug := c.Param("slug")

	err := data.DeleteSlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func RedirectLink(c *gin.Context) {
	slug := c.Param("slug")
	URL, err := data.Redirect(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ссылка не найдена"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, URL)
}