package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chiragsoni81245/jukebox/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreateMusicianNameValidation(t *testing.T) {
    expectedResponse := `{"error":"Invalid name, it should be at least of 3 characters"}`
    router, _, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    musician := models.Musician{Name: "C", Type: "Composer"}
    json_payload, _ := json.Marshal(musician)
    req, _ := http.NewRequest("POST", "/v1/musicians/", bytes.NewBuffer(json_payload))
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


func TestCreateMusician(t *testing.T) {
    expectedResponse := `{"id":1,"message":"Musician created successfully"}`
    router, _, Terminate, err := SetupRouter()
    if err != nil {
        t.Errorf("Error in db connection: %v", err)
    }
    defer Terminate() 

    musician := models.Musician{Name: "Kumar Sanu", Type: "Vocalist"}
    json_payload, _ := json.Marshal(musician)
    req, _ := http.NewRequest("POST", "/v1/musicians/", bytes.NewBuffer(json_payload))
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

func TestUpdateMusician(t *testing.T) {
    expectedResponse := `{"message":"Musician updated successfully"}`
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

    musician := models.Musician{ID: musician_id, Name: "Kumar Sanu", Type: "Vocalist"}
    json_payload, _ := json.Marshal(musician)
    req, _ := http.NewRequest("PUT", "/v1/musicians/", bytes.NewBuffer(json_payload))
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
