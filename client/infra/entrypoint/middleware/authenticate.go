package middleware

import (
	"client/infra/entrypoint/constant"
	"client/pkg/failure"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m Middleware) Authenticate(ctx *gin.Context) {
	authedClient := m.supabaseAuth.WithToken(
		ctx.Request.Header.Get("Authorization"),
	)

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
