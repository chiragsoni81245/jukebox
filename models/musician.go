package models

import (
    "errors"
)

type Musician struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    MusicianType string `json:"musician_type"`
}

func (musician *Musician) IsValid() error {
    if len(musician.Name) < 3 {
        return errors.New("Invalid name, it should be at least of 3 characters")
    }
    return nil
}

