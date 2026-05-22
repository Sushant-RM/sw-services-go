package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/BrajK111/sw-services/docs"

	"github.com/BrajK111/sw-services/internal/config"
	"github.com/BrajK111/sw-services/internal/kafka"
	"github.com/BrajK111/sw-services/internal/repository/postgres"
	"github.com/BrajK111/sw-services/internal/service"
	"github.com/BrajK111/sw-services/internal/transport/http/handler"
	"github.com/BrajK111/sw-services/internal/transport/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SW Services API
// @version 1.0
// @description Go-converted DIGIT Sewerage Service
// @host localhost:3468
// @BasePath /
func main() {

	cfg := config.Load()

	db, err := config.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repository := postgres.NewSewerageRepository(db)

	idGenClient := service.NewIDGenClient()
	workflowClient := service.NewWorkflowClient()
	producer := kafka.NewProducer()

	sewerageService := service.NewSewerageService(
		repository,
		idGenClient,
		workflowClient,
		producer,
	)

	sewerageHandler := handler.NewSewerageHandler(
		sewerageService,
	)

	mux := http.NewServeMux()

	// HEALTH ENDPOINTS (DIGIT infra/actuator compatibility)
	healthHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		handler.HealthHandler(w, r)
	}
	mux.HandleFunc("/sw-services/actuator/health", healthHandleFunc)
	mux.HandleFunc("/health", healthHandleFunc)

	mux.Handle(
		"/swagger/",
		httpSwagger.WrapHandler,
	)

	// CREATE ENDPOINTS
	createHandler := middleware.RequireMethod(http.MethodPost, sewerageHandler.CreateConnection)
	mux.HandleFunc("/sw-services/swc/_create", createHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_create", createHandler)

	// SEARCH ENDPOINTS (support both GET and POST for maximum compatibility)
	searchHandler := sewerageHandler.SearchConnection
	mux.HandleFunc("/sw-services/swc/_search", searchHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_search", searchHandler)
	mux.HandleFunc("/sw-services/swc/_plainsearch", searchHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_plainsearch", searchHandler)

	// UPDATE ENDPOINTS
	updateHandler := middleware.RequireMethod(http.MethodPost, sewerageHandler.UpdateConnection)
	mux.HandleFunc("/sw-services/swc/_update", updateHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_update", updateHandler)

	// COUNT ENDPOINTS
	countHandler := sewerageHandler.CountConnections
	mux.HandleFunc("/sw-services/swc/_count", countHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_count", countHandler)

	// ENCRYPT OLD DATA ENDPOINTS
	encryptOldDataHandler := middleware.RequireMethod(http.MethodPost, sewerageHandler.EncryptOldData)
	mux.HandleFunc("/sw-services/swc/_encryptOldData", encryptOldDataHandler)
	mux.HandleFunc("/sw-services/sewerageConnection/_encryptOldData", encryptOldDataHandler)

	handlerWithMiddleware :=
		middleware.Logging(
			middleware.Recovery(mux),
		)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handlerWithMiddleware,
	}

	go func() {

		log.Printf(
			"Server starting on :%s",
			cfg.ServerPort,
		)

		err := server.ListenAndServe()
		if err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server exited cleanly")
}
