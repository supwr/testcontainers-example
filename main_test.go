package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/supwr/testcontainers-example/pkg/database"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var conn *gorm.DB
var testDB *database.TestDatabase

func TestMain(m *testing.M) {
	testDB = database.SetupTestDatabase()
	conn = testDB.Conn
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestCreatePlayer(t *testing.T) {
	t.Run("create player successfully", func(t *testing.T) {
		defer testDB.CleanUp()

		var payload = struct {
			Name string `json:"name"`
		}{Name: "John Doe"}

		p, _ := json.Marshal(payload)
		router := setupRouter(conn)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/players", bytes.NewReader(p))
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})
}

func TestGetPlayer(t *testing.T) {
	t.Run("list players successfully", func(t *testing.T) {
		defer testDB.CleanUp()
		seedPlayers(conn)

		var players []Player
		router := setupRouter(conn)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players", nil)
		router.ServeHTTP(w, req)

		expectedPlayers := []Player{
			{ID: 1, Name: "Player 1"},
			{ID: 2, Name: "Player 2"},
			{ID: 3, Name: "Player 3"},
			{ID: 4, Name: "Player 4"},
		}

		if err := json.NewDecoder(w.Result().Body).Decode(&players); err != nil {
			log.Fatalln(err)
		}

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expectedPlayers, players)
	})
}

func seedPlayers(db *gorm.DB) {
	players := []*Player{
		{Name: "Player 1"},
		{Name: "Player 2"},
		{Name: "Player 3"},
		{Name: "Player 4"},
	}

	if err := db.Create(players).Error; err != nil {
		log.Fatalln(err.Error())
	}
}
