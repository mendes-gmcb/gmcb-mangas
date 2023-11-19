package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterImage struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	ChapterID  uuid.UUID `gorm:"index:idx_unique_page,unique,priority:1"`
	ImagePath  string    `gorm:"index:idx_unique_image,unique"`
	ImageOrder int       `gorm:"index:idx_unique_page,unique,priority:2"`
}
