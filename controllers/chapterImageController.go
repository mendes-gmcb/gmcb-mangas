package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"trabalho/initializers"
	"trabalho/models"
	"trabalho/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ChapterImageCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterImageList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterImageGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterImageUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterImageDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// func upload(fileinfo os.FileInfo) {
// 	file
// }

func ChapterImagesCreate(files []*multipart.FileHeader, mangaID, chapterID uuid.UUID, cn chan<- bool) {
	var wg sync.WaitGroup
	err_cn := make(chan error)
	semaphore := make(chan struct{}, 50)
	var chapterImages []models.ChapterImage

	for i, file := range files {
		wg.Add(1)
		semaphore <- struct{}{}

		chapterImageID := utils.GenerateUUID()

		filepath := fmt.Sprintf("/mangas/%s/%s/%s-%s", mangaID, chapterID, chapterImageID, file.Filename)

		go utils.UploadMultipleImagesToS3(file, filepath, &wg, err_cn, semaphore)

		if <-err_cn != nil {
			continue
		}

		chapterImage := models.ChapterImage{
			ID:         chapterImageID,
			ChapterID:  chapterID,
			ImagePath:  filepath,
			ImageOrder: i,
		}

		chapterImages = append(chapterImages, chapterImage)
	}

	initializers.DB.Create(chapterImages)

	wg.Wait()

	cn <- true
}
