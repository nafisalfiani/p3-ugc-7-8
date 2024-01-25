package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Value struct {
	Auth          Auth
	Log           Log
	Server        Server
	GrpcServer    Server
	MessageBroker MessageBroker
}

type Auth struct {
	SecretKey string
}

type Log struct {
	Level string
}

type Server struct {
	Base string
	Port int
}

type MessageBroker struct {
	Url string
}

func InitEnv() (*Value, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	return &Value{
		Auth: Auth{
			SecretKey: os.Getenv("AUTH_SECRETKEY"),
		},
		Log: Log{
			Level: os.Getenv("LOG_LEVEL"),
		},
		Server: Server{
			Base: os.Getenv("SERVER_BASE"),
			Port: port,
		},
		GrpcServer: Server{
			Base: os.Getenv("GRPC_SERVER_BASE"),
			Port: grpcPort,
		},
		MessageBroker: MessageBroker{
			Url: os.Getenv("RABBITMQ_URL"),
		},
	}, nil
}
