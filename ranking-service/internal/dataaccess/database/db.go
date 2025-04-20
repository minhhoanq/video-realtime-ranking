package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/minhhoanq/video-realtime-ranking/ranking-service/config"
)

type Database struct {
	*sql.DB
}

func New(cfg config.Config) (Database, error) {
	migrator, err := NewMigrator(cfg)
	if err != nil {
		return Database{}, err
	}

	err = migrator.Up(context.Background())
	if err != nil {
		return Database{}, err
	}

	db, err := newDatabase(cfg)
	if err != nil {
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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("failed to open dsn connection", err.Error())
		return nil, err
	}

	fmt.Println("database is runing on: ", cfg.DB.Host, cfg.DB.Port)

	return db, nil
}

func (p *Database) Close() {
	p.DB.Close()
}
