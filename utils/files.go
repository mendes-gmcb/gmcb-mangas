package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveCoverImage(c *gin.Context, coverImage *multipart.FileHeader, mangaID uuid.UUID) (string, error) {
	filenameArr := strings.Split(coverImage.Filename, ".")
	fileFormat := filenameArr[len(filenameArr)-1]

	coverImagePath := fmt.Sprintf("cover-images/%s.%s", mangaID, fileFormat)

	if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {
		return "", err
	}
	return coverImagePath, nil
}

func RemoveFileServerSide(filepath string) {
	os.Remove(filepath)
}
