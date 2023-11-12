package api

import (
	"client/internal/dto"
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Api) CreateAccount(ctx *gin.Context) {
	var (
		req dto.CreateAccount
		res *presenter.Login
		err error
	)
	defer a.Exception(ctx, err)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	res, err = a.authService.CreateAccount(ctx, &req)
	if err != nil {
		return
	}

	// TODO: push send email message to broker

	a.Ok(ctx, http.StatusCreated, res)
}
