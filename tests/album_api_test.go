package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chiragsoni81245/jukebox/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreateAlbumNameValidation(t *testing.T) {
    expectedResponse := `{"error":"Invalid name, it should be at least of 5 characters"}`
    router, _, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    musician := models.Album{Name: "C", ReleaseDate: time.Now(), Price: 10}
    json_payload, _ := json.Marshal(musician)
    req, _ := http.NewRequest("POST", "/v1/albums/", bytes.NewBuffer(json_payload))
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)
    
    // Read response body
    responseData, err := io.ReadAll(w.Body)
    if err != nil {
        t.Errorf("Error reading response body: %v", err)
    }

    assert.Equal(t, expectedResponse, string(responseData))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAlbumPriceValidation(t *testing.T) {
    expectedResponse := `{"error":"Invalid price, it should be in between 100 to 1000"}`
    router, _, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    album := models.Album{Name: "Chirag", ReleaseDate: time.Now(), Price: 0}
    json_payload, _ := json.Marshal(album)
    req, _ := http.NewRequest("POST", "/v1/albums/", bytes.NewBuffer(json_payload))
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)
    
    // Read response body
    responseData, err := io.ReadAll(w.Body)
    if err != nil {
        t.Errorf("Error reading response body: %v", err)
    }

    assert.Equal(t, expectedResponse, string(responseData))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}


func TestCreateAlbum(t *testing.T) {
    expectedResponse := `{"id":1,"message":"Album created successfully!"}`
    router, db, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    var musician_id uint
    err = db.QueryRow(`INSERT INTO musicians(name, type) VALUES("Chirag", "Composer") RETURNING id;`).Scan(&musician_id)
    if err != nil {
        t.Errorf("Error in inserting the dummy entry for musician in musicians table in db")
    }

    album := models.Album{Name: "Love Story", ReleaseDate: time.Now(), Price: 100, MusicianIds: []uint{musician_id}}
    json_payload, _ := json.Marshal(album)
    req, _ := http.NewRequest("POST", "/v1/albums/", bytes.NewBuffer(json_payload))
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)
    
    // Read response body
    responseData, err := io.ReadAll(w.Body)
    if err != nil {
        t.Errorf("Error reading response body: %v", err)
    }

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, expectedResponse, string(responseData))
}

func TestUpdateAlbum(t *testing.T) {
    expectedResponse := `{"message":"Album updated successfully!"}`
    router, db, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    var musician_id uint
    err = db.QueryRow(`INSERT INTO musicians(name, type) VALUES("Chirag", "Composer") RETURNING id;`).Scan(&musician_id)
    if err != nil {
        t.Errorf("Error in inserting the dummy entry for musician in musicians table in db")
    }
    
    var album_id uint
    err = db.QueryRow(`INSERT INTO albums(name, release_date, price) VALUES("Love Story", "2024-03-29T10:00:00Z", 100) RETURNING id;`).Scan(&album_id)
    if err != nil {
        t.Errorf("Error in inserting the dummy entry for album in albums table in db")
    }

    album := models.Album{ID: album_id, Name: "Swifty"}
    json_payload, _ := json.Marshal(album)
    req, _ := http.NewRequest("PUT", "/v1/albums/", bytes.NewBuffer(json_payload))
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)
    
    // Read response body
    responseData, err := io.ReadAll(w.Body)
    if err != nil {
        t.Errorf("Error reading response body: %v", err)
    }

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, expectedResponse, string(responseData))
}

