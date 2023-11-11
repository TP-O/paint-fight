package auth

import (
	"client/infra/persistence/pg"
	"client/internal/service/player"
	"time"
)

const (
	tokenLifetime = 24 * time.Hour
	bcryptCost    = 20
)

const (
	verifyEmailMsgTemplate   = "%s_%d_verify_email"
	resetPasswordMsgTemplate = "%s_%d_reset_password"
)

type Service struct {
	pg            *pg.Store
	secretKey     string
	playerService *player.Service
}

func NewService(pg *pg.Store, secretKey string, playerService *player.Service) *Service {
	return &Service{pg, secretKey, playerService}
}
