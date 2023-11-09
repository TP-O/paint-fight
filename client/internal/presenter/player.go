package presenter

import "github.com/jackc/pgx/v5/pgtype"

type Player struct {
	ID                pgtype.UUID      `json:"id"`
	Username          string           `json:"username"`
	Email             string           `json:"email"`
	Active            bool             `json:"active"`
	EmailVerifiedAt   pgtype.Timestamp `json:"email_verified_at"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	PasswordUpdatedAt pgtype.Timestamp `json:"password_updated_at"`
}
