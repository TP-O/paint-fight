package api

import (
	"client/internal/dto"
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Api) Login(ctx *gin.Context) {
	var (
		req dto.Login
		res *presenter.Login
		err error
	)
	defer a.Exception(ctx, err)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	res, err = a.authService.Login(ctx, &req)
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, res)
}
