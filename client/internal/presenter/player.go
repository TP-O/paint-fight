package presenter

import "github.com/jackc/pgx/v5/pgtype"

type Player struct {
	ID                pgtype.UUID      `json:"id"`
	Username          string           `json:"username"`
	PasswordUpdatedAt pgtype.Timestamp `json:"password_updated_at"`
}

type PlayersByUsername struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
}
