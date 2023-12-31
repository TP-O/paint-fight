package middleware

import (
	"hub/infra/entrypoint/constant"
	"hub/infra/entrypoint/exception"
	"hub/internal/service/player"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go"
)

type Middleware struct {
	supabaseAuth  gotrue.Client
	playerService *player.Service
}

var middleware *Middleware

func NewMiddleware(supabaseAuth gotrue.Client, playerService *player.Service) *Middleware {
	sync.OnceFunc(func() {
		middleware = &Middleware{
			supabaseAuth,
			playerService,
		}
	})()

	return middleware
}

func (m Middleware) Error(ctx *gin.Context, err error) {
	ctx.Set(constant.ErrorContextKey, err)
	exception.Handler(ctx)
}
