package middleware

import (
	"client/infra/entrypoint/constant"
	"client/pkg/failure"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
)

func (m Middleware) DecryptPasetoToken(ctx *gin.Context) {
	var payload paseto.JSONToken

	err := paseto.NewV2().Decrypt("", []byte(m.secretKey), &payload, nil)
	if err != nil {
		m.Error(ctx, &failure.AppError{
			HttpStatus:    http.StatusUnauthorized,
			Code:          failure.ErrInvalidToken,
			OriginalError: failure.ErrorWithTrace(err),
		})
		return
	}

	ctx.Set(constant.UserIdContextKey, payload.Subject)
	ctx.Next()
}
