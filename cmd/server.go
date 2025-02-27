package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
	"github.com/zsbahtiar/lionparcel-test/internal/handler"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/logger"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/middleware"
)

var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server!`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {

	db := database.NewPostgres(cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	movieRepository := repository.NewMovieRepository(db)
	userRepository := repository.NewUserRepository(db)

	backofficeUsecae := module.NewBackofficeUsecase(movieRepository)
	authUsecae := module.NewAuthUsecase(userRepository, cfg.JWTSecret)
	movieUsecase := module.NewMovieUsecase(movieRepository)

	validation := validator.New()
	backofficeHandler := handler.NewBackofficeHandler(backofficeUsecae, validation)
	authHandler := handler.NewAuthHandler(authUsecae, validation)
	movieHandler := handler.NewMovieHandler(movieUsecase)

	/*
		for backoffices
	*/
	router.HandleFunc("/api/backoffice/movie", backofficeHandler.CreateMovie).Methods(http.MethodPost)
	router.HandleFunc("/api/backoffice/movie/{id}", backofficeHandler.UpdateMovie).Methods(http.MethodPut)
	router.HandleFunc("/api/backoffice/movie/stat", backofficeHandler.GetStats).Methods(http.MethodGet)
	router.HandleFunc("/api/backoffice/movie", movieHandler.GetMovies).Methods(http.MethodGet)

	/*
		for user
	*/

	router.HandleFunc("/api/movie", movieHandler.GetMovies).Methods(http.MethodGet)

	router.HandleFunc("/api/movie/{id}/view", movieHandler.GetMovieView).Methods(http.MethodGet)

	router.HandleFunc("/api/movie/{id}/watch-duration", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)

	router.HandleFunc("/api/movie/{id}/vote", movieHandler.VoteMovie).Methods(http.MethodPost)

	router.HandleFunc("/api/auth/register", authHandler.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/logout", authHandler.Logout).Methods(http.MethodPost)

	router.Use(middleware.Setup)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: router,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		logger.Info(fmt.Sprintf("Server is running on port %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-sig
	logger.Info("Server shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server shutdown failed: %v", err))
	}

	logger.Info("Server gracefully stopped")
}
