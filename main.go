package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/supwr/testcontainers-example/pkg/database"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type Player struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func main() {
	// creating logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// creating database connection
	dbConfig, _ := database.NewConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// executing migration
	migration := database.NewMigration(db, dbConfig, logger)
	migration.CreateSchema()
	migration.Migrate()

	api := setupRouter(db)
	api.Run()
}

func setupRouter(db *gorm.DB) *gin.Engine {
	api := gin.Default()
	api.GET("/players", func(ctx *gin.Context) {
		var players []Player

		// list players from database
		if err := db.Find(&players).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, players)
		return
	})

	api.POST("/players", func(ctx *gin.Context) {
		type input struct {
			Name string `json:"name"`
		}
		var i input

		if err := ctx.BindJSON(&i); err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		player := &Player{Name: i.Name}

		// creating player
		if err := db.Create(&player).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, nil)
		return
	})

	return api
}
