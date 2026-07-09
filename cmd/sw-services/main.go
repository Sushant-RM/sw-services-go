package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/egovernments/sw-services-go/internal/config"
	"github.com/egovernments/sw-services-go/internal/repository/postgres"
	"github.com/egovernments/sw-services-go/internal/service"
	"github.com/egovernments/sw-services-go/internal/transport/http/handler"
	"github.com/egovernments/sw-services-go/internal/util"
	"github.com/egovernments/sw-services-go/internal/validator"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/application.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.Connect(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repo := postgres.NewConnectionRepository(db)

	propertyClient := util.NewPropertyClient(cfg.Peers.PropertyHost, cfg.Peers.PropertySearchPath)
	mdmsClient := util.NewMdmsClient(cfg.Peers.MdmsHost, cfg.Peers.MdmsSearchPath)
	idgenClient := util.NewIdGenClient(cfg.Peers.IdgenHost, cfg.Peers.IdgenGeneratePath)
	userClient := util.NewUserClient(cfg.Peers.UserHost, cfg.Peers.UserSearchPath, cfg.Peers.UserCreateNoValidate)

	propertyValidator := validator.NewPropertyValidator(propertyClient)
	mdmsValidator := validator.NewMDMSValidator(mdmsClient)

	svc := service.NewConnectionService(repo, cfg.Pagination, propertyValidator, mdmsValidator, idgenClient, userClient)
	connectionHandler := handler.NewConnectionHandler(svc)

	router := gin.Default()
	group := router.Group(cfg.Server.ContextPath)
	{
		group.GET("/actuator/health", handler.Health)
		group.POST("/swc/_create", connectionHandler.Create)
		group.POST("/swc/_update", connectionHandler.Update)
		group.POST("/swc/_search", connectionHandler.Search)
		group.POST("/swc/_plainsearch", connectionHandler.PlainSearch)
		group.POST("/swc/_encryptOldData", connectionHandler.EncryptOldData)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("sw-services-go listening on %s (context path %s)", addr, cfg.Server.ContextPath)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
