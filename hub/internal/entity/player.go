package entity

import (
	pggenerated "hub/infra/persistence/pg/generated"
	"hub/internal/presenter"
)

type Player struct {
	presenter.Presentation[pggenerated.Player, presenter.Player]
	*pggenerated.Player
}

func NewPlayer(model *pggenerated.Player) *Player {
	return &Player{
		presenter.NewPresenter[pggenerated.Player, presenter.Player](model),
		model,
	}
}
