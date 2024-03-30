package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func checkMusicianIds(musicianIds []uint, db *sql.DB, required bool) error {
    if len(musicianIds)==0{
        if required {
            return errors.New("At least one musician id is required")
        }else{
            return nil
        }
    }

    musician_ids_query := fmt.Sprintf(`
        SELECT 
            id 
        FROM musicians 
        WHERE id IN (%s);
    `, strings.Repeat("?,", len(musicianIds)-1)+"?")

    var params []interface{}
    for _, musician_id := range musicianIds {
        params = append(params, musician_id)
    }

    musician_rows, err := db.Query(musician_ids_query, params...)
    if err != nil {
        return err
    }

    var existing_musician_ids map[uint]bool = make(map[uint]bool);
    for musician_rows.Next() {
        var musician_id uint
        musician_rows.Scan(&musician_id)
        existing_musician_ids[musician_id] = true
    }

    for _, musician_id := range musicianIds {
        if !existing_musician_ids[musician_id] {
            return errors.New(fmt.Sprintf("Musician id %v does not exists", musician_id))
        }
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
    err := checkMusicianIds(album.MusicianIds, db, true)
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
    err := checkMusicianIds(album.MusicianIds, db, false)
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
    
    insert_album_musician_mapping_stmt, err := tx.Prepare("INSERT INTO album_musician (album_id, musician_id) VALUES (?, ?)")
    if err != nil {
        tx.Rollback()
        return err
    }
    defer insert_album_musician_mapping_stmt.Close()

    for _, musician_id := range album.MusicianIds {
        _, err := insert_album_musician_mapping_stmt.Exec(album.ID, musician_id)
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
    err = tx.QueryRow(update_album_query, params...).Scan(
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
    
    if len(album.MusicianIds)>0 {
        // Delete Musician Ids which are not in current array
        delete_musician_ids_query := fmt.Sprintf(`
        DELETE FROM album_musician
        WHERE 
        album_id = (?) AND    
        musician_id NOT IN (%s);
        `, strings.Repeat("?,", len(album.MusicianIds)-1)+"?" )

        var delete_query_params []interface{} 
        delete_query_params = append(delete_query_params, album.ID)
        for _, musician_id := range album.MusicianIds {
            delete_query_params = append(delete_query_params, musician_id)
        }

        _, err = tx.Exec(delete_musician_ids_query, delete_query_params...)
        if err != nil {
            tx.Rollback()
            return err
        }

        // Upsert New Musician Ids
        update_album_musician_mapping_stmt, err := tx.Prepare("INSERT OR IGNORE INTO album_musician (album_id, musician_id) VALUES (?, ?);")
        if err != nil {
            tx.Rollback()
            return err
        }
        defer update_album_musician_mapping_stmt.Close()

        for _, musician_id := range album.MusicianIds {
            _, err := update_album_musician_mapping_stmt.Exec(album.ID, musician_id)
            if err != nil {
                tx.Rollback()
                return err
            }
        }
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

