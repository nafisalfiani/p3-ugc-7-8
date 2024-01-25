package domain

import (
	"github.com/go-redis/redis/v8"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"
	"github.com/streadway/amqp"

	"github.com/sirupsen/logrus"
)

type Domains struct {
	User UserInterface
}

func Init(logger *logrus.Logger, userClient grpc.UserServiceClient, cache *redis.Client, broker *amqp.Connection) *Domains {
	return &Domains{
		User: initUser(logger, userClient, cache, broker),
	}
}
