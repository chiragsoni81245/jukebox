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
}

func (album *Album) IsValid() error {
    if len(album.Name) < 5 {
        return errors.New("Invalid name, it should be at least of 5 characters")
    }
    if album.ReleaseDate.IsZero() {
        return errors.New("Release date is mandatory")
    }
    if !(album.Price >= 100 && album.Price <= 1000) {
        return errors.New("Invalid price, it should be in between 100 to 1000")
    }
    return nil
}

func (album *Album) IsValidForUpdate() error {
    if len(album.Name)!=0 && len(album.Name) < 5 {
        return errors.New("Invalid name, it should be at least of 5 characters")
    }
    if album.Price!=0 && !(album.Price >= 100 && album.Price <= 1000) {
        return errors.New("Invalid price, it should be in between 100 to 1000")
    }
    return nil
}

func (album *Album) InsertIntoDB(db *sql.DB) error {
    query := `
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

    err := db.QueryRow(query, params...).Scan(&album.ID)  

    if err != nil {
        return err
    }

    return nil
}

func (album *Album) UpdateIntoDB(db *sql.DB) error {
    query := `
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
    err := db.QueryRow(query, params...).Scan(
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
        return err
    }

    return nil
}

func (album *Album) AddMusicians(musicianIds []uint) error {
    // TODO - write the logic to map musicians to this album
    return nil
}
