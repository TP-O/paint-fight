package api

import (
	"client/config"
	"client/internal/service/auth"
	"client/internal/service/player"
	"client/pkg/failure"
	"client/pkg/validate"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Api struct {
	cfg           config.App
	authService   *auth.Service
	playerService *player.Service
}

var api *Api

func New(
	cfg config.App,
	authService *auth.Service,
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
			authService,
			playerService,
		}
	})()

	return api
}

func (a Api) UseRouter(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.POST("/login", a.Login)
	authGroup.POST("/register", a.CreateAccount)
	// authGroup.GET("/password/forgot", a.ForgotPassword)
	// authGroup.POST("/password/reset", a.ResetPassword)
	// authGroup.GET("/email/verify", a.RequestVerifyEmail)
	// authGroup.POST("/email/verify", a.VerifyEmail)

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
			"ok":    false,
			"code":  appErr.Code,
			"error": appErr.Error(),
		})
		return
	}

	var validationErr *validator.ValidationErrors
	if errors.As(err, &validationErr) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ok":          false,
			"fieldErrors": validate.FormatValidationError(*validationErr),
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"ok":    false,
		"code":  failure.ErrUnknownCode,
		"error": err.Error(),
	})
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
