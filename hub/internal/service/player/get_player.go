package player

import (
	"context"
	"errors"
	pggenerated "hub/infra/persistence/pg/generated"
	"hub/internal/entity"
	"hub/internal/presenter"
	"hub/pkg/failure"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s Service) GetByID(ctx context.Context, id pgtype.UUID) (*entity.Player, error) {
	player, err := s.pg.PlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &failure.AppError{
				HttpStatus:    http.StatusNotFound,
				Code:          failure.ErrPlayerDoesNotExist,
				OriginalError: failure.ErrorWithTrace(err),
			}
		} else {
			return nil, &failure.AppError{
				OriginalError: failure.ErrorWithTrace(err),
			}
		}
	}

	return entity.NewPlayer(&player), nil
}

func (s Service) GetByUsername(ctx context.Context, username string) ([]presenter.PlayersByUsername, error) {
	players, err := s.pg.PlayersByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &failure.AppError{
				HttpStatus:    http.StatusNotFound,
				Code:          failure.ErrPlayerDoesNotExist,
				OriginalError: failure.ErrorWithTrace(err),
			}
		} else {
			return nil, &failure.AppError{
				OriginalError: failure.ErrorWithTrace(err),
			}
		}
	}

	res, err := presenter.PresenterFrom[[]pggenerated.PlayersByUsernameRow, []presenter.PlayersByUsername](&players)
	if err != nil {
		return nil, &failure.AppError{
			HttpStatus:    http.StatusInternalServerError,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return *res, nil
}
