// controller.go
package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"trabalho/db/models"
	"trabalho/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MangaController struct {
	Service service.MangaService
}

func (ctrl *MangaController) MangaCreate(c *gin.Context) {
	body, coverImage, err := parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manga, err := ctrl.Service.CreateManga(c, body, coverImage)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"manga": manga})
}

func (ctrl *MangaController) MangaList(c *gin.Context) {
	offset, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		offset = 0
	}

	pageSize, err := strconv.Atoi(c.Param("pageSize"))
	if defaultPageSize := 100; err != nil {
		pageSize = defaultPageSize
	}

	mangas, err := ctrl.Service.ListMangas(offset, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"mangas": mangas})
}

// não entendi a implementação
// func (ctrl *MangaController) MangaListDeleted(c *gin.Context) {
// 	offset, err := strconv.Atoi(c.Param("page"))
// 	if err != nil {
// 		offset = 0
// 	}

// 	c.JSON(http.StatusOK, gin.H{"mangas": mangas})
// }

func (ctrl *MangaController) MangaGet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	manga, err := ctrl.Service.GetMangaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manga": manga})
}

func (ctrl *MangaController) MangaUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	var updatedManga models.Manga
	if err := c.BindJSON(&updatedManga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga data provided"})
		return
	}

	manga, err := ctrl.Service.UpdateManga(id, updatedManga)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update manga"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manga": manga})
}

func (ctrl *MangaController) MangaUpdateImage(c *gin.Context) {
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

	updatedManga, err := ctrl.Service.UpdateMangaImage(c, id, coverImage)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"manga": updatedManga})
}

func (ctrl *MangaController) MangaDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manga ID"})
		return
	}

	if err := ctrl.Service.DeleteManga(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete manga"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manga deleted successfully"})
}

func parseRequest(c *gin.Context) (models.Manga, *multipart.FileHeader, error) {
	var body models.Manga
	c.Bind(&body)
	coverImage, err := c.FormFile("Image_cover")
	if err != nil {
		return models.Manga{}, nil, err
	}
	return body, coverImage, nil
}
