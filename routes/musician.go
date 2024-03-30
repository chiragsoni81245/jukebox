package routes

import (
	"github.com/chiragsoni81245/jukebox/handlers"
	"github.com/gin-gonic/gin"
)


func AttachMusicianRoutes(router *gin.RouterGroup) {
    router.POST("/", handlers.CreateMusician)
    router.PUT("/", handlers.UpdateMusician)
    router.GET("/:id/albums", handlers.GetMusicianAlbums)
}

