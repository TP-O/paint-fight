package main

import (
	"client/config"
	"client/infra/entrypoint/api"
	"client/infra/persistence/pg"
	"client/internal/service"
	"client/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
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

	cfg := config.Load("./config")

	var pgDb *pg.Store
	logger.StartToEnd(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			pgDb = pg.New(ctx, cfg.PostgreSQL)
		},
		"Connecting to PostgreSQL...",
		"Connected to PostgreSQL...",
	)

	router := gin.Default()
	apiGroup := router.Group("/")

	apiServer := api.New(
		cfg.App,
		service.NewPlayerService(pgDb),
	)
	apiServer.UseRouter(apiGroup)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Server is listening on port %d", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	<-ctx.Done()

	logger.StartToEnd(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Println(err.Error())
			}
		},
		"Shutting down...",
		"Exited",
	)
}
