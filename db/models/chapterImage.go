package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterImage struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	ChapterID  uuid.UUID
	ImagePath  string
	ImageOrder int
}
