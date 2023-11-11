package api

import (
	"client/config"
	"client/internal/service/player"
	"client/pkg/failure"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api struct {
	cfg           config.App
	playerService *player.Service
}

func New(
	cfg config.App,
	playerService *player.Service,
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
	// authGroup := router.Group("/auth")

	playerGroup := router.Group("/player")
	playerGroup.GET("/:id", a.GetPlayerByID)
	playerGroup.GET("/username/:emailOrUsername", a.GetPlayerByEmailOrUsername)
}

// TODO: Check this function after a successful handling at handler
func (a Api) Exception(ctx *gin.Context, err error) {
	var appErr *failure.AppError
	if errors.As(err, &appErr) {
		if appErr.HttpStatus == 0 {
			appErr.HttpStatus = http.StatusInternalServerError
		}

		if appErr.Code == 0 {
			appErr.Code = failure.ErrUnknownCode
		}

		ctx.JSON(appErr.HttpStatus, gin.H{
			"ok":      false,
			"code":    appErr.Code,
			"message": appErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"ok":      false,
		"code":    failure.ErrUnknownCode,
		"message": err.Error(),
	})
}

func (a Api) Ok(ctx *gin.Context, httpStatus int, data any) {
	ctx.JSON(httpStatus, gin.H{
		"ok":   true,
		"data": data,
	})
}
