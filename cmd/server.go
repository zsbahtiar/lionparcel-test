package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/handler"
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
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	backofficeUsecae := module.NewBackofficeUsecase()
	backofficeHandler := handler.NewBackofficeHandler(backofficeUsecae)

	/*
		for backoffices
	*/
	router.HandleFunc("/api/backoffices/movies", backofficeHandler.CreateMovie).Methods(http.MethodPost)
	router.HandleFunc("/api/backoffices/movies/{id}", backofficeHandler.UpdateMovie).Methods(http.MethodPut)
	router.HandleFunc("/api/backoffices/stats/most-viewed", backofficeHandler.GetMostViewed).Methods(http.MethodGet)
	router.HandleFunc("/api/backoffices/stats/most-viewed-genre", backofficeHandler.GetMostViewedGenre).Methods(http.MethodGet)
	router.HandleFunc("/api/backoffices/stats/most-voted", backofficeHandler.GetMostVoted).Methods(http.MethodGet)

	/*
		for user
	*/

	router.HandleFunc("/api/movie", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/movie/{id}/view", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)

	router.HandleFunc("/api/movie/{id}/watch-duration", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)

	router.HandleFunc("/api/movie/{id}/vote", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)
	router.HandleFunc("/api/movie/{id}/voted-movie", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Println("Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-sig
	log.Println("Server shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
