package middleware

import (
	"client/infra/entrypoint/constant"
	"client/infra/entrypoint/exception"
	"client/internal/service/player"
	"sync"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	secretKey     string
	playerService *player.Service
}

var middleware *Middleware

func NewMiddleware(secretKey string, playerService *player.Service) *Middleware {
	sync.OnceFunc(func() {
		middleware = &Middleware{
			secretKey,
			playerService,
		}
	})()

	return middleware
}

func (m Middleware) Error(ctx *gin.Context, err error) {
	ctx.Set(constant.ErrorContextKey, err)
	exception.Handler(ctx)
}
