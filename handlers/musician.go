package handlers

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/chiragsoni81245/jukebox/models"
	"github.com/chiragsoni81245/jukebox/utils"
	"github.com/gin-gonic/gin"
)

func GetMusicianAlbums(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var musician_id int
    musician_id, err = strconv.Atoi(c.Param("id"))
    if err!=nil {
        c.JSON(404, gin.H{"error": "Invalid musician id"})
        return 
    }

    get_musician_albums_query := `
        SELECT
            id, name, release_date, price, genre, description
        FROM albums as a
        INNER JOIN album_musician as am ON am.album_id=a.id
        WHERE am.musician_id=(?)
        ORDER BY a.price asc;
    `

    album_rows, err := db.Query(get_musician_albums_query, musician_id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(400, gin.H{"error": "Invalid musician id"})
        }else{
            log.Fatal(err)
            c.JSON(500, gin.H{"error": "Something went wrong"})
        }
        return 
    }

    get_album_musicians_query := `
        SELECT
            id
        FROM musicians as m
        INNER JOIN album_musician as am ON am.musician_id=m.id
        WHERE am.album_id=(?);
    ` 
    get_album_musician_stmt, err := db.Prepare(get_album_musicians_query)
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{"error": "Something went wrong"})
        return 
    }

    var albums []models.Album = make([]models.Album, 0)
    for album_rows.Next() {
        album := models.Album{}
        album_rows.Scan(
            &album.ID,
            &album.Name,
            &album.ReleaseDate,
            &album.Price,
            &album.Genre,
            &album.Description,
        )

        musician_rows, _ := get_album_musician_stmt.Query(album.ID)
        for musician_rows.Next() {
            musician := models.Musician{}
            musician_rows.Scan(&musician.ID)
            album.MusicianIds = append(album.MusicianIds, musician.ID)
        }
        albums = append(albums, album)
    }

    c.JSON(200, gin.H{
        "data": albums,
    })
}

func CreateMusician(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var musician models.Musician
    if err := c.BindJSON(&musician); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    if err := musician.IsValid(); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return 
    }
    
    err = musician.InsertIntoDB(db)

    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return 
    }

    c.JSON(200, gin.H{
        "message": "Musician created successfully",
        "id": musician.ID,
    })
}

func UpdateMusician(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var musician models.Musician
    if err := c.BindJSON(&musician); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    if musician.ID == 0 {
        c.JSON(400, gin.H{"error": "Invalid musician Id"})
        return 
    }

    if err := musician.IsValidForUpdate(); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return 
    }

    err = musician.UpdateIntoDB(db)

    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(400, gin.H{"error": "Invalid musician id"})
        }else{
            log.Fatal(err)
            c.JSON(500, gin.H{"error": "Something went wrong"})
        }
        return 
    }

    c.JSON(200, gin.H{
        "message": "Musician updated successfully",
    })
}

