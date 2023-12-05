// service.go
package service

import (
	"trabalho/db/models"
	"trabalho/repository"
	"trabalho/utils"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MangaService struct {
	Repository repository.MangaRepository
}

func (s *MangaService) CreateManga(ctx *gin.Context, body models.Manga, coverImage *multipart.FileHeader) (models.Manga, error) {
	mangaID := utils.GenerateUUID()

	coverImagePath, err := utils.SaveCoverImage(ctx, coverImage, mangaID)
	if err != nil {
		return models.Manga{}, err
	}
	defer utils.RemoveFileServerSide(coverImagePath)

	manga := createManga(body, coverImagePath, mangaID)
	if err := s.Repository.SaveManga(&manga); err != nil {
		return models.Manga{}, err
	}

	if err := utils.UploadCoverImageToS3(coverImagePath); err != nil {
		return models.Manga{}, err
	}

	return manga, nil
}


func (s *MangaService) UpdateManga(id uuid.UUID, updatedManga models.Manga) (models.Manga, error) {
	manga, err := s.Repository.GetMangaByID(id)
	if err != nil {
		return models.Manga{}, err
	}

	manga.Author = updatedManga.Author
	manga.Title = updatedManga.Title
	manga.Description = updatedManga.Description
	manga.Synopsis = updatedManga.Synopsis
	manga.Tags = updatedManga.Tags
	
	if err := s.Repository.SaveManga(&manga); err != nil {
		return models.Manga{}, err
	}
	
	return manga, nil
}

func (s *MangaService) UpdateMangaImage(ctx *gin.Context, id uuid.UUID, coverImage *multipart.FileHeader) (models.Manga, error) {
	manga, err := s.Repository.GetMangaByID(id)
	if err != nil {
		return models.Manga{}, err
	}

	oldCoverImagePath := manga.ImageCover

	newCoverImagePath, err := utils.SaveCoverImage(ctx, coverImage, id)
	if err != nil {
		return models.Manga{}, err
	}
	defer utils.RemoveFileServerSide(newCoverImagePath)

	manga.ImageCover = newCoverImagePath
	
	if err := s.Repository.UpdateManga(&manga); err != nil {
		return models.Manga{}, err
	}
	
	if err := utils.DeleteCoverImageFromS3(oldCoverImagePath); err != nil {
		return models.Manga{},err
	}
	
	return manga, nil
}

func (s *MangaService) GetMangaByID(id uuid.UUID) (models.Manga, error) {
	return s.Repository.GetMangaByID(id)
}
func (s *MangaService) DeleteManga(id uuid.UUID) error {
	manga, err := s.Repository.GetMangaByID(id)
	if err != nil {
		return err
	}

	return s.Repository.DeleteManga(&manga)
}

func (s *MangaService) ListMangas(offset, limit int) ([]models.Manga, error) {
	mangas, err := s.Repository.ListMangas(offset, limit)
	if err != nil {
		return nil, err
	}

	return mangas, nil
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
