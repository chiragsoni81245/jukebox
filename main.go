// main.go
package main

import (
    "os"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load(".env")

    r := gin.Default()

    r.GET("/health-check", func (c *gin.Context){
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    var port string = os.Getenv("PORT")
    var host string = os.Getenv("HOST")
    r.Run(fmt.Sprintf("%s:%s", host, port))
}

