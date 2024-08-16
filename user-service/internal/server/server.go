package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/milicns/company-manager/user-service/internal/api"
	"github.com/milicns/company-manager/user-service/internal/application"
	"github.com/milicns/company-manager/user-service/internal/persistance"
	"github.com/milicns/company-manager/user-service/internal/utils"
)

func registerRoutes(handler *api.Handler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.Create).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/authorize", handler.Authorize).Methods("GET")
	return router
}

func StartServer() {
	appConfig := utils.LoadAppConfig()
	dbConfig := utils.LoadDbConfig()

	db, closeDbConn := persistance.GetConn(dbConfig)
	defer closeDbConn()
	store := persistance.NewUserStore(persistance.GetCollection(db))
	service := application.NewService(store)
	authenticator := application.NewAuthenticator(store)
	handler := api.NewHandler(service, authenticator)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.AppPort),
		Handler: registerRoutes(handler),
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
