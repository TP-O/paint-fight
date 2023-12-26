package player

import (
	"context"
	pggenerated "hub/infra/persistence/pg/generated"
	"hub/internal/dto"
	"hub/internal/entity"
	"hub/pkg/failure"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s Service) Create(ctx context.Context, payload *dto.CreatePlayer) (*entity.Player, error) {
	var id pgtype.UUID
	id.Scan(payload.UserID)

	player, err := s.pg.CreatePlayer(ctx, pggenerated.CreatePlayerParams{
		UserID:   id,
		Username: payload.Username,
	})
	if err != nil {
		pgErr := failure.PgError(err)
		if pgErr != nil && pgErr.Code == failure.ErrPgUniqueViolation {
			if pgErr.ColumnName == "user_id" {
				return nil, &failure.AppError{
					HttpStatus: http.StatusBadRequest,
					Code:       failure.ErrPlayerAlreadyExists,
				}
			}

			if pgErr.ColumnName == "username" {
				return nil, &failure.AppError{
					HttpStatus: http.StatusBadRequest,
					Code:       failure.ErrUsernameAlreadyExists,
				}
			}
		}

		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreatePlayer,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return entity.NewPlayer(&player), nil
}
