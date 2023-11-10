package service

import (
	"client/infra/persistence/pg"
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/dto"
	"client/internal/presenter"
	"client/pkg/failure"
	"client/pkg/hmac"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLifetime = 24 * time.Hour
	bcryptCost    = 20
)

const (
	verifyEmailMsgTemplate   = "%s_%d_verify_email"
	resetPasswordMsgTemplate = "%s_%d_reset_password"
)

type VerifyEmailArgs struct {
	userID    pgtype.UUID
	expiredAt int64
	signature string
}

type ResetPasswordArgs struct {
	userID    pgtype.UUID
	expiredAt int64
	signature string
}

type AuthService struct {
	pg        *pg.Store
	secretKey string
}

func NewAuthService(pg *pg.Store, secretKey string) *AuthService {
	return &AuthService{pg, secretKey}
}

func (a AuthService) Login(ctx context.Context, payload *dto.Login) (*presenter.Login, error) {
	var (
		player pggenerated.Player
		res    presenter.Login
		err    error
	)

	if payload.Email != "" {
		player, err = a.pg.PlayerByEmailOrUsername(ctx, payload.Email)
	} else {
		player, err = a.pg.PlayerByEmailOrUsername(ctx, payload.Username)
	}
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrIncorrectAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(payload.Password))
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrIncorrectAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	now := time.Now()
	exp := now.Add(tokenLifetime)
	jsonToken := paseto.JSONToken{
		Issuer: "login",
		// Jti:        player.ID.String(), // TODO: implement jti checking
		Subject:    string(player.ID.Bytes[:]),
		IssuedAt:   now,
		Expiration: exp,
	}

	res.Token, err = paseto.NewV2().Encrypt([]byte(a.secretKey), jsonToken, nil)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoLoginResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	res.Player, err = presenter.PresenterFrom[pggenerated.Player, presenter.Player](&player)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoLoginResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return &res, nil
}

func (a AuthService) CreateAccount(ctx context.Context, payload *dto.CreateAccount) (*presenter.Login, error) {
	var (
		player pggenerated.Player
		res    presenter.Login
		err    error
	)

	_, err = a.pg.PlayerByEmailOrUsername(ctx, payload.Email)
	if err == nil {
		return nil, &failure.AppError{
			Code: failure.ErrEmailAlreadyExists,
		}
	} else if !errors.Is(pgx.ErrNoRows, err) {
		return nil, &failure.AppError{
			Code:          failure.ErrUnknownCode,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcryptCost)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	player, err = a.pg.CreatePlayer(ctx, pggenerated.CreatePlayerParams{
		Username: "random",
		Email:    payload.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateAccount,
			OriginalError: failure.ErrorWithTrace(err),
		}
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

	res.Token, err = paseto.NewV2().Encrypt([]byte(a.secretKey), jsonToken, nil)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoCreateAccountResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	res.Player, err = presenter.PresenterFrom[pggenerated.Player, presenter.Player](&player)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrNoCreateAccountResponse,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return &res, nil
}

func (a AuthService) VerifyEmailArgs(userID pgtype.UUID) (*VerifyEmailArgs, error) {
	var (
		args VerifyEmailArgs
		err  error
	)

	args.userID = userID
	args.expiredAt = time.Now().Add(1 * time.Hour).UnixMilli()
	args.signature, err = hmac.Sign(
		a.secretKey,
		fmt.Sprintf(verifyEmailMsgTemplate, string(args.userID.Bytes[:]), args.expiredAt),
	)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateVerifyEmailLink,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return &args, nil
}

func (a AuthService) VerifyEmail(ctx context.Context, userID pgtype.UUID, expiredAt int64, signature string) error {
	var (
		expectedSignature string
		err               error
	)

	expectedSignature, err = hmac.Sign(
		a.secretKey,
		fmt.Sprintf(verifyEmailMsgTemplate, string(userID.Bytes[:]), expiredAt),
	)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrInvalidSignature,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	if expectedSignature != signature {
		return &failure.AppError{
			Code: failure.ErrInvalidSignature,
		}
	}

	if time.Now().UnixMilli() > expiredAt {
		return &failure.AppError{
			Code: failure.ErrExpiredLink,
		}
	}

	err = a.pg.VerifyEmail(ctx, userID)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrUnableToVerifyEmail,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return nil
}

func (a AuthService) ResetPasswordArgs(ctx context.Context, email string) (*ResetPasswordArgs, error) {
	var (
		player pggenerated.Player
		args   ResetPasswordArgs
		err    error
	)

	player, err = a.pg.PlayerByEmailOrUsername(ctx, email)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, &failure.AppError{
				Code: failure.ErrEmailDoesNotExist,
			}
		} else {
			return nil, &failure.AppError{
				Code:          failure.ErrUnknownCode,
				OriginalError: failure.ErrorWithTrace(err),
			}
		}
	}

	args.userID = player.ID
	args.expiredAt = time.Now().Add(1 * time.Hour).UnixMilli()
	args.signature, err = hmac.Sign(
		a.secretKey,
		fmt.Sprintf(resetPasswordMsgTemplate, string(args.userID.Bytes[:]), args.expiredAt),
	)
	if err != nil {
		return nil, &failure.AppError{
			Code:          failure.ErrUnableToCreateResetPasswordLink,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return &args, nil
}

func (a AuthService) ResetPassword(
	ctx context.Context,
	userID pgtype.UUID,
	expiredAt int64,
	signature string,
	payload *dto.ResetPassword,
) error {
	var (
		expectedSignature string
		hashedPassword    []byte
		err               error
	)

	expectedSignature, err = hmac.Sign(
		a.secretKey,
		fmt.Sprintf(resetPasswordMsgTemplate, string(userID.Bytes[:]), expiredAt),
	)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrInvalidSignature,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	if expectedSignature != signature {
		return &failure.AppError{
			Code: failure.ErrInvalidSignature,
		}
	}

	if time.Now().UnixMilli() > expiredAt {
		return &failure.AppError{
			Code: failure.ErrExpiredLink,
		}
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(payload.Password), bcryptCost)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrUnableToUpdatePassword,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	err = a.pg.UpdatePassword(ctx, pggenerated.UpdatePasswordParams{
		ID:       userID,
		Password: string(hashedPassword),
	})
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrUnableToUpdatePassword,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return nil
}
