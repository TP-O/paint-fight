package util

import (
	"client/pkg/failure"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetCtx[T any](ctx *gin.Context, key string) (T, error) {
	var t T
	val, ok := ctx.Get(key)
	if !ok {
		return t, &failure.AppError{
			OriginalError: errors.New("context key does not exist"),
		}
	}

	t, ok = val.(T)
	if !ok {
		return t, &failure.AppError{
			OriginalError: errors.New("context value is not compatible with expected type"),
		}
	}

	return t, nil
}
