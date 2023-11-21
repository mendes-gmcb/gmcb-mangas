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

type ImageCreateRequest struct {
	MangaID    string
	ChapterID  string
	ImageOrder int
}

type ImageCreateRequestParseUUID struct {
	MangaID    uuid.UUID
	ChapterID  uuid.UUID
	ImageOrder int
}

func ChapterImageCreate(c *gin.Context) {
	var bodyS ImageCreateRequest
	var body ImageCreateRequestParseUUID
	c.Bind(&bodyS)

	file, err := c.FormFile("image")
	if err != nil {
		return
	}

	chapterImageID := utils.GenerateUUID()

	mangaID, err := uuid.Parse(bodyS.MangaID)
	if err != nil {
		return
	}

	chapterID, err := uuid.Parse(bodyS.ChapterID)
	if err != nil {
		return
	}

	body.MangaID = mangaID
	body.ChapterID = chapterID
	body.ImageOrder = bodyS.ImageOrder

	// upload image
	filepath := fmt.Sprintf("/mangas/%s/%s/%s-%s", body.MangaID, body.ChapterID, chapterImageID, file.Filename)

	fmt.Println("file info")

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)
	wg.Add(1)
	semaphore <- struct{}{}

	go utils.UploadMultipleImagesToS3(file, filepath, &wg, semaphore)
	fmt.Println("upload file")

	// save image
	chapterImage := models.ChapterImage{
		ID:         chapterImageID,
		ChapterID:  body.ChapterID,
		ImagePath:  filepath,
		ImageOrder: body.ImageOrder,
	}

	tx := initializers.DB.Begin()

	// Reorder images to add a new image
	var existingImages []models.ChapterImage
	result := initializers.DB.
		Where("chapter_id = ? AND image_order >= ?", body.ChapterID, body.ImageOrder).
		Order("image_order desc").
		Find(&existingImages)
	fmt.Println("get images to reorder")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chapter order"})
		return
	}

	// Update the image orders
	fmt.Println("init reorder")

	for i := range existingImages {
		existingImages[i].ImageOrder++
		initializers.DB.Save(&existingImages[i])
	}
	fmt.Println("finish reorder")

	// Save the new image
	result = initializers.DB.Create(&chapterImage)
	if result.Error != nil {
		wg.Add(1)
		semaphore <- struct{}{}

		go utils.DeleteFileFromS3(filepath, &wg, semaphore)
		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new image"})
		return
	}
	fmt.Println("create image")

	fmt.Println("init wait")
	wg.Wait()
	fmt.Println("finish wait")

	c.JSON(http.StatusOK, gin.H{"message": "Chapter order updated and create new image successfully"})
}

func ChapterImageUpdate(c *gin.Context) {
	// get file info
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)
	wg.Add(1)
	semaphore <- struct{}{}

	var image models.ChapterImage
	initializers.DB.Find(&image, id)

	// Delete image from S3
	if err := utils.DeleteCoverImageFromS3(image.ImagePath); err != nil {
		fmt.Printf("Error deleting image from S3: %v\n", err)
	}

	utils.UploadMultipleImagesToS3(file, image.ImagePath, &wg, semaphore)

	wg.Wait()

	// upload image
	c.JSON(http.StatusOK, gin.H{"message": "Image updated on s3 successfully"})
}

func ChapterImageDelete(c *gin.Context) {
	// get info
	// reorder images to remove the image
	// remove image from s3

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ChapterImagesCreate(files []*multipart.FileHeader, mangaID, chapterID uuid.UUID) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)
	var chapterImages []models.ChapterImage

	for i, file := range files {
		wg.Add(1)
		semaphore <- struct{}{}

		chapterImageID := utils.GenerateUUID()

		filepath := fmt.Sprintf("/mangas/%s/%s/%s-%s", mangaID, chapterID, chapterImageID, file.Filename)

		go utils.UploadMultipleImagesToS3(file, filepath, &wg, semaphore)

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
}
