package database

import (
	"context"
	"database/sql"
	"fmt"
	"video-realtime-ranking/config"
)

type Database struct {
	*sql.DB
}

func New(cfg config.Config) (Database, error) {
	migrator, err := NewMigrator(cfg)
	if err != nil {
		// l.Error("failed to new migration database", zap.Error(err))
		return Database{}, err
	}

	err = migrator.Up(context.Background())
	if err != nil {
		// l.Error("failed to migration up schema database", zap.Error(err))
		return Database{}, err
	}

	db, err := newDatabase(cfg)
	if err != nil {
		// l.Error("failed to connection to database", zap.Error(err))
		return Database{}, err
	}

	return Database{
		DB: db,
	}, nil
}

func newDatabase(cfg config.Config) (*sql.DB, error) {
	// create data source name (DSN) string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	// Open GORM database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// l.Error("failed to open dsn connection", zap.Error(err))
		return nil, err
	}

	// l.Info("database is running on",
	// 	zap.String("Host: ", cfg.DBHost),
	// 	zap.String("Name: ", cfg.DBName),
	// 	zap.Int("Port: ", cfg.DBPort))

	return db, nil
}

func (p *Database) Close() {
	p.DB.Close()
}
