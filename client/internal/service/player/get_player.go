package player

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/entity"
	"client/pkg/failure"
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Player(ctx context.Context, id pgtype.UUID) (*entity.Player, error) {
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

func (s Service) PlayerByEmailOrUsername(ctx context.Context, emailOrUsername string) (*entity.Player, error) {
	player, err := s.pg.PlayerByEmailOrUsername(ctx, emailOrUsername)
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

func (s Service) CreateAccount(ctx context.Context, email, password string) (*entity.Player, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	player, err := s.pg.CreatePlayer(ctx, pggenerated.CreatePlayerParams{
		Username: "random", // TODO: random username
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return entity.NewPlayer(&player), nil
}
