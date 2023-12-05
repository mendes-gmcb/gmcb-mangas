package main

import (
	"trabalho/controllers"
	"trabalho/db"
	"trabalho/env"

	"github.com/gin-gonic/gin"
)

func init() {
	env.Load()
	db.Load()
}

func main() {
	r := gin.Default()

	r.MaxMultipartMemory = 256 << 20

	r.GET("/manga/page/:page", controllers.MangaList)
	r.POST("/manga", controllers.MangaCreate)
	r.GET("/manga/:id", controllers.MangaGet)
	r.PUT("/manga/:id", controllers.MangaUpdate)
	r.PATCH("/manga/:id", controllers.MangaUpdateImage)
	r.DELETE("/manga/:id", controllers.MangaDelete)

	r.GET("/chapter/page/:page/manga/:mangaid", controllers.ChapterList)
	r.POST("/chapter", controllers.ChapterCreate)
	r.GET("/chapter/:id", controllers.ChapterGet)
	r.PUT("/chapter/:id", controllers.ChapterUpdateOrder)
	r.DELETE("/chapter/:id", controllers.ChapterDelete)

	r.POST("/chapter-image", controllers.ChapterImageCreate)
	r.PATCH("/chapter-image/:id", controllers.ChapterImageUpdate)
	r.DELETE("/chapter-image/:id", controllers.ChapterImageDelete)

	r.Run()
}
