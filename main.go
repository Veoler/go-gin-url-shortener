package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Veoler/go-gin-url-shortener/data"
	"github.com/Veoler/go-gin-url-shortener/models"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, linksSlug)
}

func UpdateLink (c *gin.Context) {
	id := c.Param("id")

	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при изменении"})
		return
	}
	updateLink, err := data.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updateLink)
}

func AddLink (c *gin.Context) {
	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при добавлении: некорректные данные"})
		return
	}
	c.JSON(http.StatusCreated, data.Add(input))
}

func DeleteLinkByID(c *gin.Context) {
	id := c.Param("id")

	err := data.DeleteID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Delete успешен"})
}

func DeleteLinkBySlug(c *gin.Context) {
	slug := c.Param("slug")

	err := data.DeleteSlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Delete успешен"})
}

func RedirectLink(c *gin.Context) {
	slug := c.Param("slug")
	err := data.Redirect(slug)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": "Статья не найдена"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "https://github.com/intocode/go-gin-url-shortener")
}