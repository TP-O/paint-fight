package exception

import (
	"errors"
	"hub/infra/entrypoint/constant"
	"hub/pkg/failure"
	"hub/pkg/validate"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Handler(ctx *gin.Context) {
	ctxErr, ok := ctx.Get(constant.ErrorContextKey)
	if !ok {
		return
	}

	var err error
	if err, ok = ctxErr.(error); !ok {
		return
	}

	var appErr *failure.AppError
	if errors.As(err, &appErr) {
		if appErr.HttpStatus == 0 {
			appErr.HttpStatus = http.StatusInternalServerError
		}

		if appErr.Code == 0 {
			appErr.Code = failure.ErrUnknownCode
		}

		ctx.AbortWithStatusJSON(appErr.HttpStatus, gin.H{
			"ok":    false,
			"code":  appErr.Code,
			"error": appErr.Error(),
		})
		return
	}

	var validationErr *validator.ValidationErrors
	if errors.As(err, &validationErr) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok":          false,
			"fieldErrors": validate.FormatValidationError(*validationErr),
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"ok":    false,
		"code":  failure.ErrUnknownCode,
		"error": err.Error(),
	})
}
