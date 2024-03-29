package routes

import (
	"github.com/chiragsoni81245/jukebox/handlers"
	"github.com/gin-gonic/gin"
)


func AttachAlbumRoutes(router *gin.RouterGroup) {
    router.POST("/", handlers.CreateAlbum)
    router.PUT("/", handlers.UpdateAlbum)
}

