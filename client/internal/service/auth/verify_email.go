package auth

import (
	"client/pkg/failure"
	"client/pkg/hmac"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailArgs struct {
	userID    pgtype.UUID
	expiredAt int64
	signature string
}

func (s Service) VerifyEmailArgs(userID pgtype.UUID) (*VerifyEmailArgs, error) {
	var (
		args VerifyEmailArgs
		err  error
	)

	args.userID = userID
	args.expiredAt = time.Now().Add(1 * time.Hour).UnixMilli()
	args.signature, err = hmac.Sign(
		s.secretKey,
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

func (s Service) VerifyEmail(ctx context.Context, userID pgtype.UUID, expiredAt int64, signature string) error {
	var (
		expectedSignature string
		err               error
	)

	expectedSignature, err = hmac.Sign(
		s.secretKey,
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

	err = s.pg.VerifyEmail(ctx, userID)
	if err != nil {
		return &failure.AppError{
			Code:          failure.ErrUnableToVerifyEmail,
			OriginalError: failure.ErrorWithTrace(err),
		}
	}

	return nil
}
