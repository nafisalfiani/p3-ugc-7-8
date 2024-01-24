package main

import (
	"api-gateway/config"
	"api-gateway/docs"
	"api-gateway/domain"
	"api-gateway/handler"
	"api-gateway/usecase"
	"fmt"
	"log"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"

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

	// init domain
	dom := domain.Init(logger, userClient)

	// init usecase
	usecase := usecase.Init(cfg, logger, dom)

	// init handler
	handler := handler.Init(cfg, usecase, validator, logger)

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

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", cfg.Server.Base, cfg.Server.Port)))
}
