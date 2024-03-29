package main

import (
	"fmt"
	"log"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/config"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/domain"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/usecase"
	"github.com/streadway/amqp"
)

// @contact.name Nafisa Alfiani
// @contact.email nafisa.alfiani.ica@gmail.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// init config
	cfg, err := config.InitEnv()
	if err != nil {
		log.Fatalln(err)
	}

	// init logger
	logger, err := config.InitLogger(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info(fmt.Sprintf("%#v", cfg))

	// init DB connection
	db, err := config.InitNoSql(cfg)
	if err != nil {
		logger.Fatalf("failed to connect to mongo. %v", err)
	}

	// init cache
	redis, err := config.InitCache(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to connect to redis. %v", err)
	}

	// init pubsub
	rabbitMqServer, err := amqp.Dial(cfg.MessageBroker.Url)
	if err != nil {
		panic(err)
	}

	defer rabbitMqServer.Close()

	// init domain
	dom := domain.Init(db, redis, logger, rabbitMqServer)

	// init handler
	uc := usecase.Init(cfg, logger, dom)

	g := grpc.Init(cfg, logger, uc)
	g.Run()
}
