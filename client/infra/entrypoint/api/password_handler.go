package api

import (
	"client/internal/dto"
	"client/internal/service/auth"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a Api) ForgotPassword(ctx *gin.Context) {
	var (
		req  dto.ForgotPassword
		args *auth.ResetPasswordArgs
		err  error
	)
	defer a.Exception(ctx, err)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	args, err = a.authService.ResetPasswordArgs(ctx, req.Email)
	if err != nil {
		return
	}

	// TODO: push send email message to broker
	fmt.Println(args)

	a.Ok(ctx, http.StatusOK, nil)
}

func (a Api) ResetPassword(ctx *gin.Context) {
	var (
		id        [16]byte
		expiredAt int
		req       dto.ResetPassword
		err       error
	)

	copy(id[:], ctx.Query("id"))

	expiredAt, err = strconv.Atoi(ctx.Query("expiredAt"))
	if err != nil {
		return
	}

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = a.authService.ResetPassword(
		ctx,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		int64(expiredAt),
		ctx.Query("signature"),
		&req)
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, nil)
}
