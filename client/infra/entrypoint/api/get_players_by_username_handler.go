package api

import (
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Api) GetPlayersUsername(ctx *gin.Context) {
	var (
		presenter []presenter.PlayersByUsername
		err       error
	)
	defer a.Error(ctx, err)

	presenter, err = a.playerService.GetByUsername(ctx, ctx.Param("username"))
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
