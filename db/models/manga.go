package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manga struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Author      string
	Title       string
	Description string
	Synopsis    string
	Tags        string
	ImageCover  string

	Chapters []Chapter `gorm:"constraint:OnDelete:CASCADE; foreignKey:MangaID"`
}
