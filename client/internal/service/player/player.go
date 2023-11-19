package player

import (
	"client/infra/persistence/pg"
)

type Service struct {
	pg *pg.Store
}

func NewService(pg *pg.Store) *Service {
	return &Service{pg}
}
