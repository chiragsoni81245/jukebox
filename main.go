// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/chiragsoni81245/jukebox/middlewares"
	"github.com/chiragsoni81245/jukebox/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func SetupRouter(db_name string) (*gin.Engine, *sql.DB) {
    router := gin.Default()

    // Database Setup
    db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db", db_name))
    if err != nil {
        log.Fatal(err)
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

    return router, db
}

func main() {
    godotenv.Load(".env")
    router, db := SetupRouter("jukebox")
    defer db.Close()

    var port string = os.Getenv("PORT")
    if len(port)==0 {
        port = "8080"
    }
    var host string = os.Getenv("HOST")
    if len(host)==0 {
        host = "localhost"
    }
    router.Run(fmt.Sprintf("%s:%s", host, port))
}

