package main

import (
	"client/config"
	"client/infra/entrypoint/api"
	"client/infra/entrypoint/middleware"
	"client/infra/persistence/pg"
	"client/internal/service/auth"
	"client/internal/service/player"
	"client/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(1)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfgDir := os.Getenv("CONFIG_DIR")
	if cfgDir == "" {
		cfgDir = "./config"
	}
	cfg := config.Load(cfgDir)

	var pgDb *pg.Store
	logger.StartToEnd(
		logger.StartToEndConfig{
			Action: func(ctx context.Context) {
				pgDb = pg.New(ctx, cfg.PostgreSQL)
			},
			Timeout:  30 * time.Second,
			StartMsg: "Connecting to PostgreSQL...",
			EndMsg:   "Connected to PostgreSQL",
		},
	)

	router := gin.Default()
	apiGroup := router.Group("/")

	playerService := player.NewService(pgDb)
	authService := auth.NewService(pgDb, cfg.Secret.Auth, playerService)
	middleware := middleware.NewMiddleware(cfg.Secret.Auth, playerService)

	apiServer := api.New(
		cfg.App,
		middleware,
		authService,
		playerService,
	)
	apiServer.UseRouter(apiGroup)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Server is listening on port %d", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server is closed: %s", err.Error())
		}
	}()

	<-ctx.Done()

	logger.StartToEnd(
		logger.StartToEndConfig{
			Action: func(ctx context.Context) {
				if err := server.Shutdown(ctx); err != nil {
					log.Fatalf("Shutting down failed: %s", err.Error())
				}
			},
			Timeout:  5 * time.Second,
			StartMsg: "Shutting down...",
			EndMsg:   "Exited gracefully",
		},
	)
}
