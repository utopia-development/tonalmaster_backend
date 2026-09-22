package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/utopia-development/tonalmaster_backend/internal/config"
	"github.com/utopia-development/tonalmaster_backend/internal/database"
	"github.com/utopia-development/tonalmaster_backend/internal/domain/calendars"
	"github.com/utopia-development/tonalmaster_backend/internal/handlers"
	"github.com/utopia-development/tonalmaster_backend/internal/repository"
	"github.com/utopia-development/tonalmaster_backend/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil { logger.Error("configuration failed", "error", err); os.Exit(1) }

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil { logger.Error("database initialization failed", "error", err); os.Exit(1) }
	defer db.Close()

	registry := calendars.NewRegistry(calendars.NewTonalpohualliCASO())
	health := handlers.NewHealthHandler(db)
	calendarHandler := handlers.NewCalendarHandler(registry)
	authRepository := repository.NewPostgresAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService, cfg.Env != "development", cfg.RegistrationCode)
	contentRepository := repository.NewPostgresContentRepository(db)
	contentService := services.NewContentService(contentRepository)
	contentHandler := handlers.NewContentHandler(contentService)
	uleRepository := repository.NewPostgresULERepository(db)
	uleService := services.NewULEService(uleRepository)
	uleHandler := handlers.NewULEHandler(uleService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Health)
	mux.HandleFunc("GET /ready", health.Ready)
	mux.HandleFunc("GET /api/v1/calendars", calendarHandler.List)
	mux.HandleFunc("GET /api/v1/calendars/{id}", calendarHandler.Get)
	mux.HandleFunc("GET /api/v1/calendars/convert", calendarHandler.Convert)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/me", authHandler.Me)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.Handle("POST /api/v1/events", handlers.RequireAuth(authService, http.HandlerFunc(contentHandler.CreateEvent)))
	mux.Handle("GET /api/v1/events", handlers.RequireAuth(authService, http.HandlerFunc(contentHandler.ListEvents)))
	mux.Handle("DELETE /api/v1/events/{id}", handlers.RequireAuth(authService, http.HandlerFunc(contentHandler.DeleteEvent)))
	mux.Handle("POST /api/v1/interpretations", handlers.RequireAuth(authService, http.HandlerFunc(contentHandler.CreateInterpretation)))
	mux.HandleFunc("GET /api/v1/interpretations", contentHandler.ListInterpretations)
	mux.HandleFunc("GET /api/v1/articles", uleHandler.Articles)
	mux.HandleFunc("GET /api/v1/articles/{id}", uleHandler.Article)
	mux.HandleFunc("GET /api/v1/bibliography", uleHandler.Bibliography)
	mux.HandleFunc("GET /api/v1/bibliography/{id}", uleHandler.Biblio)
	mux.HandleFunc("GET /api/v1/catalogs", uleHandler.Catalogs)
	mux.HandleFunc("GET /api/v1/catalogs/{id}", uleHandler.Catalog)
	mux.HandleFunc("GET /api/v1/ads", uleHandler.Ads)

	server := &http.Server{Addr: cfg.Address(), Handler: corsMiddleware(cfg.CORSAllowedOrigins, mux), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}

	go func() {
		logger.Info("server started", "addr", cfg.Address(), "env", cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Error("server stopped unexpectedly", "error", err); stop() }
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil { logger.Error("graceful shutdown failed", "error", err); os.Exit(1) }
	logger.Info("server stopped")
}

func corsMiddleware(allowedOrigins []string, next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins { allowed[origin] = struct{}{} }
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowed[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		}
		if r.Method == http.MethodOptions && origin != "" { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w, r)
	})
}
