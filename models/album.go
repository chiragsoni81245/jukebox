package models

import (
	"database/sql"
	"errors"
	"time"
	"github.com/chiragsoni81245/jukebox/utils"
)

type Album struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    ReleaseDate time.Time  `json:"release_date"`
    Genre       string  `json:"genre"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    MusicianIds []uint  `json:"musicianIds"`
}

func checkMusicianIds(musicianIds []uint, db *sql.DB) error {
    check_musician_ids_query := `
        SELECT 
            COUNT(*) = COUNT(DISTINCT id) AS all_ids_present 
        FROM musicians 
        WHERE id IN (?);
    `

    var isAllMusicianIdsCorrect bool
    err := db.QueryRow(check_musician_ids_query, musicianIds).Scan(&isAllMusicianIdsCorrect) 
    if err != nil {
        return err
    }

    if !isAllMusicianIdsCorrect {
        return errors.New("Invalid musicianIds")
    }

    return nil
}

func (album *Album) IsValid(db *sql.DB) error {
    if len(album.Name) < 5 {
        return errors.New("Invalid name, it should be at least of 5 characters")
    }
    if album.ReleaseDate.IsZero() {
        return errors.New("Release date is mandatory")
    }
    if !(album.Price >= 100 && album.Price <= 1000) {
        return errors.New("Invalid price, it should be in between 100 to 1000")
    }

    // Check if all the musician ids are present in database
    err := checkMusicianIds(album.MusicianIds, db)
    if err != nil {
        return err
    }

    return nil
}

func (album *Album) IsValidForUpdate(db *sql.DB) error {
    if len(album.Name)!=0 && len(album.Name) < 5 {
        return errors.New("Invalid name, it should be at least of 5 characters")
    }
    if album.Price!=0 && !(album.Price >= 100 && album.Price <= 1000) {
        return errors.New("Invalid price, it should be in between 100 to 1000")
    }

    // Check if all the musician ids are present in database
    err := checkMusicianIds(album.MusicianIds, db)
    if err != nil {
        return err
    }

    return nil
}

func (album *Album) InsertIntoDB(db *sql.DB) error {
    insert_album_query := `
        INSERT INTO albums(
            name, 
            release_date, 
            price, 
            genre, 
            description
        ) VALUES($1, $2, $3, $4, $5) RETURNING id;
    `

    var params []any
    params = utils.AppendIfNotEmpty(params, album.Name)
    params = utils.AppendIfNonZero(params, album.ReleaseDate)
    params = utils.AppendIfNonZero(params, album.Price)
    params = utils.AppendIfNotEmpty(params, album.Genre)
    params = utils.AppendIfNotEmpty(params, album.Description)

    tx, err := db.Begin()
    err = tx.QueryRow(insert_album_query, params...).Scan(&album.ID)  
    if err != nil {
        tx.Rollback()
        return err
    }
    
    insert_album_musician_mapping_stmt, err := tx.Prepare("INSERT INTO album_musicians (album_id, musician_id) VALUES (?, ?)")
    if err != nil {
        tx.Rollback()
        panic(err)
    }
    defer insert_album_musician_mapping_stmt.Close()

    for _, id := range album.MusicianIds {
        _, err := insert_album_musician_mapping_stmt.Exec(album.ID, id)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

func (album *Album) UpdateIntoDB(db *sql.DB) error {
    update_album_query := `
        UPDATE albums 
        SET
            name=coalesce(?, name), 
            release_date=coalesce(?, release_date), 
            price=coalesce(?, price), 
            genre=coalesce(?, genre), 
            description=coalesce(?, description) 
        WHERE id=?
        RETURNING 
            id, name, release_date, price, genre, description;
    `

    var params []interface{}
    params = utils.AppendIfNotEmpty(params, album.Name)
    params = utils.AppendIfNonZero(params, album.ReleaseDate)
    params = utils.AppendIfNonZero(params, album.Price)
    params = utils.AppendIfNotEmpty(params, album.Genre)
    params = utils.AppendIfNotEmpty(params, album.Description)
    params = append(params, album.ID)

    var genre sql.NullString
    var description sql.NullString
    tx, err := db.Begin()
    err = db.QueryRow(update_album_query, params...).Scan(
        &album.ID,
        &album.Name,
        &album.ReleaseDate,
        &album.Price,
        &genre,
        &description,
    )

    if genre.Valid {
        album.Genre = genre.String
    }

    if description.Valid {
        album.Description = description.String
    }
    if err != nil {
        tx.Rollback()
        return err
    }
    
    update_album_musician_mapping_stmt, err := tx.Prepare("INSERT OR IGNORE INTO album_musicians (album_id, musician_id) VALUES (?, ?);")
    if err != nil {
        tx.Rollback()
        panic(err)
    }
    defer update_album_musician_mapping_stmt.Close()

    for _, id := range album.MusicianIds {
        _, err := update_album_musician_mapping_stmt.Exec(album.ID, id)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

