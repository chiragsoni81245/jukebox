package handlers

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/chiragsoni81245/jukebox/models"
	"github.com/chiragsoni81245/jukebox/utils"
	"github.com/gin-gonic/gin"
)

func GetAlbums(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    get_musician_albums_query := `
        SELECT
            id, name, release_date, price, genre, description
        FROM albums as a
        ORDER BY a.release_date asc;
    `

    album_rows, err := db.Query(get_musician_albums_query)
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{"error": "Something went wrong"})
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

func GetAlbumMusicians(c *gin.Context) {
    db, err := utils.GetDB(c) 
    if err != nil {
        log.Fatal(err)
        c.JSON(500, gin.H{
            "error": "Something went wrong",
        })
        return
    }

    var album_id int
    album_id, err = strconv.Atoi(c.Param("id"))
    if err!=nil {
        c.JSON(404, gin.H{"error": "Invalid album id"})
        return 
    }

    get_album_musicians_query := `
        SELECT
            id, name, type
        FROM musicians as m
        INNER JOIN album_musician as am ON am.musician_id=m.id
        WHERE am.album_id=(?)
        ORDER BY m.name;
    ` 

    musicians := make([]models.Musician, 0)
    musician_rows, _ := db.Query(get_album_musicians_query, album_id)
    for musician_rows.Next() {
        musician := models.Musician{}
        musician_rows.Scan(
            &musician.ID,
            &musician.Name,
            &musician.Type,
        )

        musicians = append(musicians, musician)
    }

    c.JSON(200, gin.H{
        "data": musicians,
    })
}

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

