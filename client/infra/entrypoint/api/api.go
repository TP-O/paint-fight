package api

import (
	"client/config"
	"client/internal/service"

	"github.com/gin-gonic/gin"
)

type api struct {
	cfg           config.App
	playerService *service.PlayerService
}

func New(
	cfg config.App,
	playerService *service.PlayerService,
) *api {
	if cfg.Env == config.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	return &api{
		cfg,
		playerService,
	}
}

func (a api) UseRouter(router *gin.RouterGroup) {
	router.GET("/player/:id", a.Player)
}
