package db

import (
	"trabalho/db/models"
)

func Migrate() {
	DB.AutoMigrate(
		&models.Manga{},
		&models.Chapter{},
		&models.ChapterImage{},
	)
}