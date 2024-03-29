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

func main() {
    godotenv.Load(".env")
    router := gin.Default()

    // Database Setup
    db, err := sql.Open("sqlite3", "jukebox.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Providing this database instance to all the requests into there context it self
    // which they can access using `c.db` where `c` is *gin.Context
    router.Use(middlewares.DBWrapper(db))    

    v1 := router.Group("/v1")
    {
        albumRouter := v1.Group("/album")
        routes.AttachAlbumRoutes(albumRouter)
    }
    

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

