package api

import (
	"client/config"
	"client/internal/service"
	"client/pkg/failure"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api struct {
	cfg           config.App
	playerService *service.PlayerService
}

func New(
	cfg config.App,
	playerService *service.PlayerService,
) *Api {
	if cfg.Env == config.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Api{
		cfg,
		playerService,
	}
}

func (a Api) UseRouter(router *gin.RouterGroup) {
	router.GET("/player/:id", a.GetPlayerByID)
}

func (a Api) Exception(ctx *gin.Context, err error) {
	var appErr *failure.AppError
	if errors.As(err, &appErr) {
		if appErr.Status == 0 {
			appErr.Status = http.StatusInternalServerError
		}

		ctx.JSON(appErr.Status, gin.H{
			"ok":      false,
			"message": appErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"ok":      false,
		"message": err.Error(),
	})
}

func (a Api) Ok(ctx *gin.Context, statusCode int, data any) {
	ctx.JSON(statusCode, gin.H{
		"ok":   true,
		"data": data,
	})
}
