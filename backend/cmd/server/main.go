package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/meetia/backend/internal/api"
	"github.com/meetia/backend/internal/config"
	"github.com/meetia/backend/internal/db"
	"github.com/meetia/backend/internal/repository"
	"github.com/meetia/backend/internal/services/auth"
	"github.com/meetia/backend/internal/services/meeting"
	"github.com/meetia/backend/internal/services/webrtc"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	cfg := config.Load()
	database := db.Initialize(cfg)
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Heartbeat("/health"))

	humaConfig := huma.DefaultConfig("Meetia API", "1.0.0")
	humaConfig.Info.Description = "API for Meetia video conferencing platform"
	humaapi := humachi.New(r, humaConfig)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	meetingRepo := repository.NewMeetingRepository(database)

	authService := auth.NewAuthService(userRepo, cfg.JWTSecret, 24*time.Hour)
	meetingService := meeting.NewMeetingService(meetingRepo, userRepo)
	sfuService := webrtc.NewSFUService()

	api.SetupRoutes(humaapi, authService, sfuService, meetingService)

	// start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out... forcing exit.")
			}
		}()

		// trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// run the server
	log.Printf("Server is running on port %s\n", cfg.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
