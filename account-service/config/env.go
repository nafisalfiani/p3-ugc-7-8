package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Value struct {
	NoSqlDatabase NoSqlDatabase
	Auth          Auth
	Log           Log
	Server        Server
}

type NoSqlDatabase struct {
	DSN         string
	MaxIdleTime string
	MaxIdleConn string
}

type Auth struct {
	SecretKey string
}

type Server struct {
	Base string
	Port int
}

type Log struct {
	Level string
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

	return &Value{
		NoSqlDatabase: NoSqlDatabase{
			DSN:         os.Getenv("MONGO_DSN"),
			MaxIdleTime: os.Getenv("MONGO_MAX_IDLE_TIME"),
			MaxIdleConn: os.Getenv("MONGO_MAX_IDLE_CONN"),
		},
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
	}, nil
}
