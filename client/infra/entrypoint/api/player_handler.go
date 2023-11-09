package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a api) Player(ctx *gin.Context) {
	var id [16]byte
	copy(id[:], ctx.Param("id"))

	player, err := a.playerService.FindPlayer(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return
	}

	presenter, err := player.Presenter.Presenter()
	if err != nil {
		return
	}

	ctx.JSON(200, map[string]any{
		"ok":   true,
		"data": presenter,
	})
}
