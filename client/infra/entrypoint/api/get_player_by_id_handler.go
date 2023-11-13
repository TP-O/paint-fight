package api

import (
	"client/internal/entity"
	"client/internal/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a Api) GetPlayerByID(ctx *gin.Context) {
	var (
		id        [16]byte
		player    *entity.Player
		presenter *presenter.Player
		err       error
	)
	defer a.Error(ctx, err)

	copy(id[:], ctx.Param("id"))
	player, err = a.playerService.Player(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return
	}

	presenter, err = player.Presenter()
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
