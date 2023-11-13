package api

import (
	"client/infra/entrypoint/constant"
	"client/internal/service/auth"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a Api) RequestVerifyEmail(ctx *gin.Context) {
	var (
		id   [16]byte
		args *auth.VerifyEmailArgs
		err  error
	)
	defer a.Error(ctx, err)

	copy(id[:], []byte(ctx.GetString(constant.UserIdContextKey)))

	args, err = a.authService.VerifyEmailArgs(pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return
	}

	// TODO: push send email message to broker
	fmt.Println(args)

	a.Ok(ctx, http.StatusOK, nil)
}

func (a Api) VerifyEmail(ctx *gin.Context) {
	var (
		id        [16]byte
		expiredAt int
		err       error
	)
	defer a.Error(ctx, err)

	copy(id[:], ctx.Query("id"))

	expiredAt, err = strconv.Atoi(ctx.Query("expiredAt"))
	if err != nil {
		return
	}

	err = a.authService.VerifyEmail(
		ctx,
		pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		int64(expiredAt),
		ctx.Query("signature"))
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, nil)
}
