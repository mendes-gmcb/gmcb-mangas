package main

import (
	"trabalho/initializers"
	model "trabalho/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConectToDB()
}

func main() {
	initializers.DB.AutoMigrate(
		&model.Manga{},
		&model.Chapter{},
		&model.ChapterImage{},
	)
}
