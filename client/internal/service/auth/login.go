package auth

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/dto"
	"client/internal/presenter"
	"client/pkg/failure"
	"context"
	"net/http"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(ctx context.Context, payload *dto.Login) (*presenter.Login, error) {
	var (
		player pggenerated.Player
		res    presenter.Login
		err    error
	)

	if payload.Email != "" {
		player, err = s.pg.PlayerByEmailOrUsername(ctx, payload.Email)
	} else {
		player, err = s.pg.PlayerByEmailOrUsername(ctx, payload.Username)
	}
	if err != nil {
		return nil, &failure.AppError{
			HttpStatus:    http.StatusBadRequest,
			Code:          failure.ErrIncorrectAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(payload.Password))
	if err != nil {
		return nil, &failure.AppError{
			HttpStatus:    http.StatusBadRequest,
			Code:          failure.ErrIncorrectAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	// TODO: remove this duplicate logic (generate token) in login and create account feature
	now := time.Now()
	exp := now.Add(tokenLifetime)
	jsonToken := paseto.JSONToken{
		Issuer: "login",
		// Jti:        player.ID.String(), // TODO: implement jti checking
		Subject:    string(player.ID.Bytes[:]),
		IssuedAt:   now,
		Expiration: exp,
	}

	res.AccessToken.ExpiredAt = exp.UnixMilli()
	res.AccessToken.Value, err = paseto.NewV2().Encrypt([]byte(s.secretKey), jsonToken, nil)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoLoginResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	if playerPresenter, err := presenter.PresenterFrom[pggenerated.Player, presenter.Player](&player); err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoLoginResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	} else {
		res.Player = *playerPresenter
	}

	return &res, nil
}
