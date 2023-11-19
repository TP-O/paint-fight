package player

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/dto"
	"client/internal/entity"
	"client/pkg/failure"
	"context"
)

func (s Service) Create(ctx context.Context, payload *dto.CreatePlayer) (*entity.Player, error) {
	player, err := s.pg.CreatePlayer(ctx, pggenerated.CreatePlayerParams{
		Username: payload.Username,
	})
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreatePlayer,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return entity.NewPlayer(&player), nil
}
