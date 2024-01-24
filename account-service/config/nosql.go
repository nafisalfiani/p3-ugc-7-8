package config

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitNoSql(cfg *Value) (*mongo.Client, error) {
	maxPoolSize, err := strconv.ParseUint(cfg.NoSqlDatabase.MaxIdleConn, 10, 64)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.NoSqlDatabase.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	opts := options.Client().
		ApplyURI(cfg.NoSqlDatabase.DSN).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMaxPoolSize(maxPoolSize).
		SetMinPoolSize(1).
		SetMaxConnIdleTime(duration)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cleanup := context.WithTimeout(context.Background(), 10*time.Second)
	defer cleanup()

	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		return nil, err
	}

	return client, nil
}
