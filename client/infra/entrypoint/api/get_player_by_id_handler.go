package api

import (
	"client/internal/domain/entity"
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a Api) GetPlayerByID(ctx *gin.Context) {
	var (
		player    *entity.Player
		presenter *presenter.Player
		err       error
	)
	defer a.Exception(ctx, err)

	var id [16]byte
	copy(id[:], ctx.Param("id"))

	player, err = a.playerService.Player(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return
	}

	presenter, err = player.Presenter.Presenter()
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
