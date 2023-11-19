package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	MangaID       uuid.UUID `gorm:"index:idx_unique_chapter,unique,priority:1"`
	ChapterNumber int       `gorm:"index:idx_unique_chapter,unique,priority:2"`
	NumPages      int

	Images []ChapterImage `gorm:"constraint:OnDelete:CASCADE; foreignKey:ChapterID"`
}
