package api

import (
	"hub/internal/entity"
	"hub/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a Api) GetPlayerByID(ctx *gin.Context) {
	var (
		id        pgtype.UUID
		player    *entity.Player
		presenter *presenter.Player
		err       error
	)
	defer a.Error(ctx, err)

	err = id.Scan(ctx.Param("id"))
	if err != nil {
		return
	}

	player, err = a.playerService.GetByID(ctx, id)
	if err != nil {
		return
	}

	presenter, err = player.Presenter()
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
