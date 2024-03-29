// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

