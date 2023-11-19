package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"trabalho/initializers"
	"trabalho/models"
	"trabalho/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ChapterCreate(c *gin.Context) {
	cic_cn := make(chan bool)
	body, files := parseChapterRequest(c)

	chapterID := utils.GenerateUUID()

	go ChapterImagesCreate(files, body.MangaID, chapterID, cic_cn)

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

	<-cic_cn

	initializers.DB.Preload("Images").First(&chapter, chapterID)

	c.JSON(http.StatusCreated, gin.H{"chapter": chapter})
}

func ChapterList(c *gin.Context) {
	offset, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		offset = 0
	}

	mangaID, err := uuid.Parse(c.Param("mangaid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	var chapters []models.Chapter

	initializers.DB.Where("manga_id", mangaID).Order("chapter_number asc").Limit(100).Offset(offset).Find(&chapters)

	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
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

func parseChapterRequest(c *gin.Context) (body models.Chapter, files []*multipart.FileHeader) {
	form, _ := c.MultipartForm()

	files = form.File["Image[]"]

	mangaID, err := uuid.Parse(form.Value["MangaID"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	chapterNumber, err := strconv.Atoi(form.Value["ChapterNumber"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter number"})
		return
	}

	numPages, err := strconv.Atoi(form.Value["NumPages"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter number"})
		return
	}

	body.MangaID = mangaID
	body.ChapterNumber = chapterNumber
	body.NumPages = numPages
	return
}
