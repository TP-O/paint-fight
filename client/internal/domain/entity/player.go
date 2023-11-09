package entity

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/presenter"
)

type Player struct {
	presenter.Presenter[pggenerated.Player, presenter.Player]
	model *pggenerated.Player
}

func NewPlayer(model *pggenerated.Player) *Player {
	return &Player{
		presenter.NewPresenter[pggenerated.Player, presenter.Player](model),
		model,
	}
}
