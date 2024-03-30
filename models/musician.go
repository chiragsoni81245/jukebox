package models

import (
	"database/sql"
	"errors"
	"github.com/chiragsoni81245/jukebox/utils"
)

type Musician struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    Type string `json:"type"`
}

func (musician *Musician) IsValid() error {
    if len(musician.Name) < 3 {
        return errors.New("Invalid name, it should be at least of 3 characters")
    }
    return nil
}

func (musician *Musician) IsValidForUpdate() error {
    if len(musician.Name)!=0 && len(musician.Name) < 3 {
        return errors.New("Invalid name, it should be at least of 3 characters")
    }

    return nil
}

func (musician *Musician) InsertIntoDB(db *sql.DB) error {
    insert_musician_query := `
        INSERT INTO musicians(
            name, 
            type
        ) VALUES($1, $2) RETURNING id;
    `

    var params []any
    params = utils.AppendIfNotEmpty(params, musician.Name)
    params = utils.AppendIfNotEmpty(params, musician.Type)

    tx, err := db.Begin()
    err = tx.QueryRow(insert_musician_query, params...).Scan(&musician.ID)  
    if err != nil {
        tx.Rollback()
        return err
    }
    
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

func (musician *Musician) UpdateIntoDB(db *sql.DB) error {
    update_musician_query := `
        UPDATE musicians 
        SET
            name=coalesce(?, name), 
            type=coalesce(?, type) 
        WHERE id=(?)
        RETURNING 
            id, name, type;
    `

    var params []interface{}
    params = utils.AppendIfNotEmpty(params, musician.Name)
    params = utils.AppendIfNotEmpty(params, musician.Type)
    params = append(params, musician.ID)

    tx, err := db.Begin()
    err = tx.QueryRow(update_musician_query, params...).Scan(
        &musician.ID,
        &musician.Name,
        &musician.Type,
    )

    if err != nil {
        tx.Rollback()
        return err
    }
    
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

