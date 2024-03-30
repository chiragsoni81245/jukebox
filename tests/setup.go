package tests

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	"github.com/chiragsoni81245/jukebox/middlewares"
	"github.com/chiragsoni81245/jukebox/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var DB_COUNTER int = 1;

func SetupRouter() (*gin.Engine, *sql.DB, func(), error) {
    router := gin.Default()

    // Setup Migrations

    db_name := fmt.Sprintf("test_%d.db", DB_COUNTER)
    DB_COUNTER += 1
    cmd := exec.Command("goose", "-dir", "../migrations", "sqlite3", db_name, "up")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return nil, nil, nil, err
    }
    // Database Setup
    db, err := sql.Open("sqlite3", db_name)
    if err != nil {
        return nil, nil, nil, err
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

    return router, db, func(){
        db.Close()
        os.Remove(db_name)
    }, nil
}

