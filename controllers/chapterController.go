package controllers

import (
	"net/http"
	"trabalho/initializers"
	"trabalho/models"
	"trabalho/utils"

	"github.com/gin-gonic/gin"
)

func ChapterCreate(c *gin.Context) {
	var body models.Chapter
	c.Bind(&body)

	chapterID := utils.GenerateUUID()

	chapter := models.Chapter{
		ID:            chapterID,
		MangaID:       body.MangaID,
		ChapterNumber: body.ChapterNumber,
		NumPages:      body.NumPages,
	}

	result := initializers.DB.Create(&chapter)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chapter": chapter})
}

func ChapterList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
