package middleware

import (
	"client/infra/entrypoint/constant"
	"client/pkg/failure"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m Middleware) Authenticate(ctx *gin.Context) {
	token := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		token = ctx.Query("token")
		if token == "" {
			m.Error(ctx, &failure.AppError{
				HttpStatus: http.StatusUnauthorized,
				Code:       failure.ErrInvalidToken,
			})
			return
		}
	}

	authedClient := m.supabaseAuth.WithToken(token)
	user, err := authedClient.GetUser()
	if err != nil {
		m.Error(ctx, &failure.AppError{
			HttpStatus:    http.StatusUnauthorized,
			Code:          failure.ErrInvalidToken,
			OriginalError: failure.ErrorWithTrace(err),
		})
		return
	}

	ctx.Set(constant.UserContextKey, user)
	ctx.Set(constant.UserIDContextKey, user.ID.String())
	ctx.Next()
}
