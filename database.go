package kocto

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func DBConnect(cfg DBConfig) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(cfg.URL)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Name)

	return db, nil
}
