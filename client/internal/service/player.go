package service

import (
	"client/infra/persistence/pg"
	"client/internal/domain/entity"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerService struct {
	pg *pg.Store
}

func NewPlayerService(pg *pg.Store) *PlayerService {
	return &PlayerService{pg}
}

func (p PlayerService) FindPlayer(ctx context.Context, id pgtype.UUID) (*entity.Player, error) {
	player, err := p.pg.PlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity.NewPlayer(&player), nil
}

func (p PlayerService) FindPlayerByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.Player, error) {
	player, err := p.pg.PlayerByEmailOrUsername(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	return entity.NewPlayer(&player), nil
}
