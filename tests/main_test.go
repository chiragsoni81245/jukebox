package tests

import (
	"database/sql"

	"github.com/chiragsoni81245/jukebox/middlewares"
	"github.com/chiragsoni81245/jukebox/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

func SetupRouter() (*gin.Engine, *sql.DB, error) {
    router := gin.Default()

    // Database Setup
    db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
    if err != nil {
        return nil, nil, err
    }
    // Setup Migrations
    err = goose.Up(db, "../migrations")
    if err != nil {
        return nil, nil, err
    }

    // Providing this database instance to all the requests into there context it self
    // which they can access using `c.db` where `c` is *gin.Context
    router.Use(middlewares.DBWrapper(db))    

    v1 := router.Group("/v1")
    {
        albumRouter := v1.Group("/albums")
        routes.AttachAlbumRoutes(albumRouter)
        musicianRouter := v1.Group("/musicians")
        routes.AttachMusicianRoutes(musicianRouter)
    }

    return router, db, nil
}

