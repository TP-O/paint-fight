package entity

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/presenter"
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
