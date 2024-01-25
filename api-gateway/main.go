package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/config"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/docs"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/domain"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/handler"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// init validator
	validator := validator.New(validator.WithRequiredStructEnabled())

	// init grpc procedure
	cc, err := grpc.Dial(fmt.Sprintf("%v:%v", cfg.GrpcServer.Base, cfg.GrpcServer.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer cc.Close()
	userClient := grpc.NewUserServiceClient(cc)

	// init cache
	cache, err := config.InitCache(cfg, logger)
	if err != nil {
		log.Fatalf("failed to init cache: %v", err)
	}

	// init message broker
	conn, err := config.InitBroker()
	if err != nil {
		logger.Fatalln(err)
	}

	// init domain
	dom := domain.Init(logger, userClient, cache, conn)

	// init usecase
	usecase := usecase.Init(cfg, logger, dom)

	// init handler
	handler := handler.Init(cfg, usecase, validator, logger, conn)

	// start consumer on go routine
	ch, err := handler.StartConsumer()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Debug(ch)

	// init echo instance
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(handler.MiddlewareLogging)
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	docs.SwaggerInfo.Title = "API Gateway"
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	e.GET("/ping", handler.Ping)

	api := e.Group("/api")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)

	users := api.Group("/users", handler.Authorize)
	users.GET("", handler.ListUsers)
	users.POST("", handler.CreateUser)
	users.GET("/:id", handler.GetUser)
	users.PUT("/:id", handler.UpdateUser)
	users.DELETE("/:id", handler.DeleteUser)

	// start http server on go routine
	go func() {
		if err := e.Start(fmt.Sprintf("%v:%v", cfg.Server.Base, cfg.Server.Port)); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	logger.Info("Shutting down...")

	// Close Echo HTTP server
	if err := e.Shutdown(context.Background()); err != nil {
		logger.Fatalln("Error shutting down HTTP server:", err)
	}

	// Close Rabbit Connection
	// conn.Close()
	// ch.Close()

	logger.Info("Server shutdown complete.")

}
