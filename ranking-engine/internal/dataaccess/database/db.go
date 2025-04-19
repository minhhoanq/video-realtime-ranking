package database

import (
	"context"
	"fmt"
	"log"
	"video-realtime-ranking/ranking-engine/config"
	"video-realtime-ranking/ranking-engine/pkg/constants"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var collection *mongo.Collection

type Database struct {
	*mongo.Client
	cfg config.Config
}

func New(cfg config.Config) (Database, error) {
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	otps := options.Client().ApplyURI(connString)

	client, err := mongo.Connect(otps)
	if err != nil {
		return Database{}, err
	}

	// check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("cannot connection to mongodb: ", err.Error())
		return Database{}, err
	}

	log.Println("connection to mongodb successfully", cfg.DB.Port, cfg.DB.Name)

	// Migrations
	migrator := NewMigrator(client)
	// migration collections if not exists
	migrator.EnsureCollectionsAndIndexes(cfg.DB.Name, []string{constants.RANKING_COLLECTION})

	return Database{
		Client: client,
		cfg:    cfg,
	}, nil
}

func (db Database) returnCollectionPointer(collection string) *mongo.Collection {
	return db.Client.Database(db.cfg.DB.Name).Collection(collection)
}

func (db Database) Disconnect(ctx context.Context) error {
	return db.Disconnect(ctx)
}
