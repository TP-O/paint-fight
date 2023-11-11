package player

import (
	"client/infra/persistence/pg"
)

const (
	bcryptCost = 20
)

type Service struct {
	pg *pg.Store
}

func NewService(pg *pg.Store) *Service {
	return &Service{pg}
}
