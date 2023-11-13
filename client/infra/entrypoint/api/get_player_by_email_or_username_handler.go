package api

import (
	"client/internal/entity"
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Api) GetPlayerByEmailOrUsername(ctx *gin.Context) {
	var (
		player    *entity.Player
		presenter *presenter.Player
		err       error
	)
	defer a.Error(ctx, err)

	player, err = a.playerService.PlayerByEmailOrUsername(ctx, ctx.Param("emailOrUsername"))
	if err != nil {
		return
	}

	presenter, err = player.Presenter()
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
