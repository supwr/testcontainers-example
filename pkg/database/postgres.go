package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewConnection(cfg Config) (*gorm.DB, error) {
	var err error
	var conn *gorm.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s search_path=%s",
		cfg.DatabaseHost, cfg.DatabaseUsername, cfg.DatabasePassword, cfg.DatabaseDBName, cfg.DatabasePort, cfg.DatabaseSchema)

	conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Discard,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", cfg.DatabaseSchema),
			SingularTable: false,
		},
	})

	if err != nil {
		return nil, err
	}

	db, err := conn.DB()

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
