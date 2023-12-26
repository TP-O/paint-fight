package api

import (
	"hub/infra/entrypoint/constant"
	"hub/internal/dto"
	"hub/internal/entity"
	"hub/internal/presenter"
	"hub/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	gotruetype "github.com/supabase-community/gotrue-go/types"
)

func (a Api) CreatePlayer(ctx *gin.Context) {
	var (
		user      gotruetype.User
		req       dto.CreatePlayer
		player    *entity.Player
		presenter *presenter.Player
		err       error
	)
	defer a.Error(ctx, err)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	user, err = util.GetCtx[gotruetype.User](ctx, constant.UserContextKey)
	if err != nil {
		return
	}

	req.UserID = user.ID.String()
	player, err = a.playerService.Create(ctx, &req)
	if err != nil {
		return
	}

	presenter, err = player.Presenter()
	if err != nil {
		return
	}

	a.Ok(ctx, http.StatusOK, presenter)
}
