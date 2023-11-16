package models

import (
	"gorm.io/gorm"
)

type ChapterImage struct {
	gorm.Model
	ChapterID  uint
	ImagePath  string
	ImageOrder int
	ImageName  string
}
