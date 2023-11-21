package supabase

import (
	"client/config"
	"sync"

	"github.com/supabase-community/gotrue-go"
)

type Supabase struct {
	auth gotrue.Client
}

var supabase *Supabase

func New(cfg config.Supabase) *Supabase {
	sync.OnceFunc(func() {
		supabase = &Supabase{
			gotrue.New(
				cfg.ID,
				cfg.AnonKey,
			),
		}
	})()

	return supabase
}

func (s Supabase) Auth() gotrue.Client {
	return s.auth
}
