package logger

import (
	"context"
	"log"
	"time"
)

type StartToEndConfig struct {
	Action   func(ctx context.Context)
	StartMsg string
	EndMsg   string
	Timeout  time.Duration
}

func StartToEnd(cfg StartToEndConfig) {
	ctx := context.Background()
	if cfg.Timeout.Nanoseconds() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	log.Println(cfg.StartMsg)
	cfg.Action(ctx)
	log.Println(cfg.EndMsg)
}
