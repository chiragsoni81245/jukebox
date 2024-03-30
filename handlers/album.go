package handlers

import (
	"database/sql"
	"log"
	"github.com/chiragsoni81245/jukebox/models"
	"github.com/chiragsoni81245/jukebox/utils"
	"github.com/gin-gonic/gin"
)

func CreateAlbum(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var album models.Album
    if err := c.BindJSON(&album); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    if err := album.IsValid(db); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return 
    }
    
    err = album.InsertIntoDB(db)

    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return 
    }

    c.JSON(200, gin.H{
        "message": "Album created successfully!",
        "id": album.ID,
    })
}

func UpdateAlbum(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var album models.Album
    if err := c.BindJSON(&album); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    if album.ID == 0 {
        c.JSON(500, gin.H{"error": "Invalid album Id"})
        return 
    }

    if err := album.IsValidForUpdate(db); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return 
    }

    err = album.UpdateIntoDB(db)

    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(400, gin.H{"error": "Invalid album id"})
        }else{
            log.Fatal(err)
            c.JSON(500, gin.H{"error": "Something went wrong!"})
        }
        return 
    }

    c.JSON(200, gin.H{
        "message": "Album updated successfully!",
    })
}

