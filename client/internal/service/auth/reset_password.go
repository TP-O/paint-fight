package auth

import (
	pggenerated "client/infra/persistence/pg/generated"
	"client/internal/dto"
	"client/pkg/failure"
	"client/pkg/hmac"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordArgs struct {
	userID    pgtype.UUID
	expiredAt int64
	signature string
}

func (s Service) ResetPasswordArgs(ctx context.Context, email string) (*ResetPasswordArgs, error) {
	var (
		player pggenerated.Player
		args   ResetPasswordArgs
		err    error
	)

	player, err = s.pg.PlayerByEmailOrUsername(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
		s.secretKey,
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

func (s Service) ResetPassword(ctx context.Context, userID pgtype.UUID, expiredAt int64, signature string, payload *dto.ResetPassword) error {
	var (
		expectedSignature string
		hashedPassword    []byte
		err               error
	)

	expectedSignature, err = hmac.Sign(
		s.secretKey,
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
			HttpStatus: http.StatusBadRequest,
			Code:       failure.ErrInvalidSignature,
		}
	}

	if time.Now().UnixMilli() > expiredAt {
		return &failure.AppError{
			HttpStatus: http.StatusBadRequest,
			Code:       failure.ErrExpiredLink,
		}
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(payload.Password), bcryptCost)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrUnableToUpdatePassword,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	err = s.pg.UpdatePassword(ctx, pggenerated.UpdatePasswordParams{
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
