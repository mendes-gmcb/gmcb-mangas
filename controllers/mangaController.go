package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"trabalho/initializers"
	"trabalho/models"
	"trabalho/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MangaCreate(c *gin.Context) {
	body, coverImage, err := parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mangaID := utils.GenerateUUID()

	coverImagePath, err := utils.SaveCoverImage(c, coverImage, mangaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer utils.RemoveFileServerSide(coverImagePath)

	manga := createManga(body, coverImagePath, mangaID)
	if err := saveMangaToDB(manga); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := utils.UploadCoverImageToS3(coverImagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"manga": manga})
}

func MangaList(c *gin.Context) {
	offset, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		// chama função de log
		offset = 0
	}

	var mangas []models.Manga

	initializers.DB.Order("title asc").Limit(30).Offset(offset).Find(&mangas)

	c.JSON(http.StatusOK, gin.H{"mangas": mangas})
}

func MangaGet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	var manga models.Manga

	initializers.DB.Preload("Chapters").First(&manga, id)

	c.JSON(http.StatusOK, gin.H{"manga": manga})
}

func MangaUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	var body models.Manga

	c.Bind(&body)

	var manga models.Manga
	initializers.DB.First(&manga, id)

	initializers.DB.Model(&manga).Updates(models.Manga{
		Author:      body.Author,
		Title:       body.Title,
		Description: body.Description,
		Synopsis:    body.Synopsis,
		Tags:        body.Tags,
	})

	c.JSON(http.StatusOK, gin.H{"manga": manga})
}

func MangaUpdateImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	_, coverImage, err := parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var manga models.Manga
	initializers.DB.Preload("Chapters").First(&manga, id)

	imageCover := manga.ImageCover

	// Delete cover image from S3
	if err := utils.DeleteCoverImageFromS3(imageCover); err != nil {
		fmt.Printf("Error deleting cover image from S3: %v\n", err)
		return
	}

	coverImagePath, err := utils.SaveCoverImage(c, coverImage, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer utils.RemoveFileServerSide(coverImagePath)

	if err := utils.UploadCoverImageToS3(coverImagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	initializers.DB.Model(&manga).Updates(models.Manga{ImageCover: coverImagePath})

	c.JSON(http.StatusOK, gin.H{"manga": manga})
}

func MangaDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	// Fetch manga from the database
	var manga models.Manga
	if err := initializers.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	imageCover := manga.ImageCover

	// Delete manga from the database
	if err := initializers.DB.Delete(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete manga"})
		return
	}

	// Delete cover image from S3
	if err := utils.DeleteCoverImageFromS3(imageCover); err != nil {
		fmt.Printf("Error deleting cover image from S3: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manga deleted successfully"})
}


func createManga(body models.Manga, coverImagePath string, mangaID uuid.UUID) models.Manga {
	return models.Manga{
		ID:          mangaID,
		Author:      body.Author,
		Title:       body.Title,
		Description: body.Description,
		Synopsis:    body.Synopsis,
		Tags:        body.Tags,
		ImageCover:  coverImagePath,
	}
}

func saveMangaToDB(manga models.Manga) error {
	result := initializers.DB.Create(&manga)
	return result.Error
}
