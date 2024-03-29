package models

import (
	"errors"
	"time"
)

type Album struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    ReleaseDate time.Time  `json:"release_date"`
    Genre       string  `json:"genre"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    Musicians   []*Musician
}

func (album *Album) IsValid() error {
    if len(album.Name) < 5 {
        return errors.New("Invalid name, it should be at least of 5 characters")
    }
    if album.ReleaseDate.IsZero() {
        return errors.New("Invalid release date")
    }
    if !(album.Price >= 100 && album.Price <= 1000) {
        return errors.New("Invalid price, it should be in between 100 to 1000")
    }
    return nil
}
