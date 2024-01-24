package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func InitCache(cfg *Value, logger *logrus.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	logger.Debug(res)

	return client, nil
}
