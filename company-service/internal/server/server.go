package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/milicns/company-manager/company-service/internal/api"
	"github.com/milicns/company-manager/company-service/internal/application"
	"github.com/milicns/company-manager/company-service/internal/persistance"
	"github.com/milicns/company-manager/company-service/internal/utils"

	"github.com/gorilla/mux"
)

func registerRoutes(handler *api.Handler, userSvcAddress string) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/{id}", handler.GetOne).Methods("GET")
	router.HandleFunc("/create", api.Authorize(handler.Create, userSvcAddress)).Methods("POST")
	router.HandleFunc("/{id}", api.Authorize(handler.Delete, userSvcAddress)).Methods("DELETE")
	router.HandleFunc("/{id}", api.Authorize(handler.Patch, userSvcAddress)).Methods("PATCH")
	return router
}

func StartServer() {
	appConfig := utils.LoadAppConfig()
	kafkaConfig := utils.LoadKafkaConfig()
	userSvcConfig := utils.LoadUserServiceConfig()
	dbConfig := utils.LoadDbConfig()

	kafkaAddr := fmt.Sprintf("%s:%s", kafkaConfig.KafkaHost, kafkaConfig.KafkaPort)
	userSvcAddress := fmt.Sprintf("http://%s:%s/authorize", userSvcConfig.UserServiceHost, userSvcConfig.UserServicePort)

	db, closeDbConn := persistance.GetConn(dbConfig)
	defer closeDbConn()
	store := persistance.NewCompanyStore(persistance.GetCollection(db))
	service := application.NewService(store)
	producer, closeKafkaConn := application.NewProducer(kafkaAddr)
	defer closeKafkaConn()
	handler := api.NewHandler(service, producer)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.AppPort),
		Handler: registerRoutes(handler, userSvcAddress),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("server shutting down")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Printf("failed to shutdown server gracefully: %v\n", err)
	}
}
