package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rigbyel/ad-market/internal/config"
	adcreate "github.com/rigbyel/ad-market/internal/http-server/handlers/advert/create"
	"github.com/rigbyel/ad-market/internal/http-server/handlers/feed/show"
	"github.com/rigbyel/ad-market/internal/http-server/handlers/user/login"
	"github.com/rigbyel/ad-market/internal/http-server/handlers/user/register"
	"github.com/rigbyel/ad-market/internal/http-server/middleware/cors"
	"github.com/rigbyel/ad-market/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// initializing config
	cfg := config.MustLoad()

	// initializing logger
	log := setupLogger(cfg.Env)
	log.Info("starting ad-market", slog.String("env", cfg.Env))

	// initializing sqlite storage
	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("err", err.Error()))
	}

	// intializing chi router
	router := chi.NewRouter()

	// middleware
	router.Use(cors.MiddlewareCors)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// handlers
	router.Post("/register", register.New(log, storage, cfg.JwtSecret))
	router.Post("/login", login.New(log, storage, cfg.JwtSecret, cfg.TokenTL))
	router.Post("/advert", adcreate.New(log, storage, cfg.JwtSecret))
	router.Get("/feed", show.New(log, storage, cfg.JwtSecret))

	// starting server
	log.Info("starting server", slog.String("addres", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("err", err.Error()))
	}

	log.Error("server stopped")

}

// setting up logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
