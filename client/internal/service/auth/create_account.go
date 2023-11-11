package auth

import (
	"client/internal/dto"
	"client/internal/entity"
	"client/internal/presenter"
	"client/pkg/failure"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/o1egl/paseto"
)

func (s Service) CreateAccount(ctx context.Context, payload *dto.CreateAccount) (*presenter.Login, error) {
	var (
		player *entity.Player
		res    presenter.Login
		err    error
	)

	_, err = s.pg.PlayerByEmailOrUsername(ctx, payload.Email)
	if err == nil {
		return nil, &failure.AppError{
			Code: failure.ErrEmailAlreadyExists,
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, &failure.AppError{
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	player, err = s.playerService.CreateAccount(ctx, payload.Email, payload.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	exp := now.Add(tokenLifetime)
	jsonToken := paseto.JSONToken{
		Issuer: "create_account",
		// Jti:        player.ID.String(), // TODO: implement jti checking
		Subject:    string(player.ID.Bytes[:]),
		IssuedAt:   now,
		Expiration: exp,
	}

	res.AccessToken.ExpiredAt = exp.UnixMilli()
	res.AccessToken.Value, err = paseto.NewV2().Encrypt([]byte(s.secretKey), jsonToken, nil)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoCreateAccountResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	if playerPresenter, err := player.Presenter(); err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoCreateAccountResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	} else {
		res.Player = *playerPresenter
	}

	return &res, nil
}
