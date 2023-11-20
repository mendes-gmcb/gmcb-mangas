package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"trabalho/initializers"
	"trabalho/models"
	"trabalho/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chaptersOrder struct {
	ChapterID uuid.UUID `json:"chapterID"`
	Order     int       `json:"Order"`
}

func ChapterCreate(c *gin.Context) {
	body, files := parseChapterRequest(c)

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

	ChapterImagesCreate(files, body.MangaID, chapterID)

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var chapter models.Chapter

	initializers.DB.Preload("Images").First(&chapter, id)

	c.JSON(http.StatusOK, gin.H{"chapter": chapter})
}

func ChapterUpdateOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	var chapters []chaptersOrder

	c.Bind(&chapters)
	if len(chapters) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No chapters provided"})
		return
	}

	// Begin a transaction
	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Iterate through the received chapters and update their order
	for _, chapter := range chapters {
		// fmt.Println(chapter)
		if err := tx.Model(&models.Chapter{}).
			Where("id = ? AND manga_id = ?", chapter.ChapterID, id).
			Update("chapter_number", nil).Error; err != nil {

			// Rollback the transaction on error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chapter order"})
			return
		}
	}

	// Iterate through the received chapters and update their order
	for _, chapter := range chapters {
		// fmt.Println(chapter)
		if err := tx.Model(&models.Chapter{}).
			Where("id = ? AND manga_id = ?", chapter.ChapterID, id).
			Update("chapter_number", chapter.Order).Error; err != nil {

			// Rollback the transaction on error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chapter order"})
			return
		}
	}

	// Commit the transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Chapter order updated successfully"})
}

func ChapterDelete(c *gin.Context) {
	var wg sync.WaitGroup

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var chapter models.Chapter
	if err := initializers.DB.Preload("Images").First(&chapter, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	images := chapter.Images

	// Delete manga from the database
	if err := initializers.DB.Delete(&chapter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete manga"})
		return
	}

	semaphore := make(chan struct{}, 50)

	for _, image := range images {
		wg.Add(1)
		semaphore <- struct{}{}

		go utils.DeleteFileFromS3(image.ImagePath, &wg, semaphore)
	}

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"message": "Chapter deleted successfully"})
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
