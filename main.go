package main

import (
	"trabalho/controllers"
	"trabalho/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConectToDB()
	initializers.ConectToS3()
	initializers.ConectUploaderToS3()
}

func main() {
	r := gin.Default()
	r.GET("/manga/page/:page", controllers.MangaList)
	r.POST("/manga", controllers.MangaCreate)
	r.GET("/manga/:id", controllers.MangaGet)
	r.PUT("/manga/:id", controllers.MangaUpdate)
	r.DELETE("/manga/:id", controllers.MangaDelete)

	r.GET("/chapter/page/:page", controllers.ChapterList)
	r.POST("/chapter", controllers.ChapterCreate)
	r.GET("/chapter/:id", controllers.ChapterGet)
	r.PUT("/chapter/:id", controllers.ChapterUpdate)
	r.DELETE("/chapter/:id", controllers.ChapterDelete)

	r.GET("/chapter-image/page/:page", controllers.ChapterImageList)
	r.GET("/chapter-image/:id", controllers.ChapterImageGet)
	r.PUT("/chapter-image/:id", controllers.ChapterImageUpdate)
	r.DELETE("/chapter-image/:id", controllers.ChapterImageDelete)

	r.Run()
}