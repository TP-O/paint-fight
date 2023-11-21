package api

import (
	"client/config"
	"client/infra/entrypoint/constant"
	"client/infra/entrypoint/exception"
	"client/infra/entrypoint/middleware"
	"client/internal/service/player"
	"client/pkg/validate"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Api struct {
	cfg           config.App
	middleware    *middleware.Middleware
	playerService *player.Service
}

var api *Api

// TODO: use struct as single paramater
func New(
	cfg config.App,
	middleware *middleware.Middleware,
	playerService *player.Service,
) *Api {
	sync.OnceFunc(func() {
		if cfg.Env == config.ProdEnv {
			gin.SetMode(gin.ReleaseMode)
		}

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			validate.SetUp(v)
		}

		api = &Api{
			cfg,
			middleware,
			playerService,
		}
	})()

	return api
}

func (a Api) UseRouter(router *gin.RouterGroup) {
	playerGroup := router.Group("/player")
	playerGroup.GET("/:id", a.GetPlayerByID)
	playerGroup.GET("/username/:username", a.GetPlayersUsername)
	playerGroup.POST("/", a.middleware.Authenticate, a.CreatePlayer)

	router.Use(exception.Handler)
}

func (a Api) Error(ctx *gin.Context, err error) {
	ctx.Set(constant.ErrorContextKey, err)
}

func (a Api) Ok(ctx *gin.Context, httpStatus int, data any) {
	if data != nil {
		ctx.JSON(httpStatus, gin.H{
			"ok":   true,
			"data": data,
		})
	} else {
		ctx.JSON(httpStatus, gin.H{
			"ok": true,
		})
	}
}
