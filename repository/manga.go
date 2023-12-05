// repository.go
package repository

import (
	"trabalho/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MangaRepository struct{
	db *gorm.DB
}

func (r *MangaRepository) SaveManga(manga *models.Manga) error {
	result := r.db.Create(manga)
	return result.Error
}

func (r *MangaRepository) GetMangaByID(id uuid.UUID) (models.Manga, error) {
	var manga models.Manga
	if err := r.db.First(&manga, id).Error; err != nil {
		return models.Manga{}, err
	}
	return manga, nil
}

func (r *MangaRepository) UpdateManga(manga *models.Manga) error {
	result := r.db.Save(manga)
	return result.Error
}

func (r *MangaRepository) DeleteManga(manga *models.Manga) error {
	result := r.db.Delete(manga)
	return result.Error
}

func (r *MangaRepository) ListMangas(offset, limit int) ([]models.Manga, error) {
	var mangas []models.Manga
	result := r.db.Offset(offset).Limit(limit).Find(&mangas)
	return mangas, result.Error
}
