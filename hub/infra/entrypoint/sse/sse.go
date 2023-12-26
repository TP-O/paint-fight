package sse

import (
	"hub/infra/entrypoint/constant"
	"hub/infra/entrypoint/middleware"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type Client chan string

type UserID = string

type Sse struct {
	newClientUserIDs    chan UserID
	closedClientUserIDs chan UserID
	clients             map[UserID]Client
	middleware          *middleware.Middleware
}

func New(middleware *middleware.Middleware) *Sse {
	sse := &Sse{
		make(chan UserID),
		make(chan UserID),
		make(map[UserID]Client),
		middleware,
	}

	go sse.listen()
	return sse
}

func (s Sse) listen() {
	for {
		select {
		case userID := <-s.newClientUserIDs:
			log.Printf("Client %s added. %d registered clients", userID, len(s.clients))

		case userID := <-s.closedClientUserIDs:
			close(s.clients[userID])
			delete(s.clients, userID)
			log.Printf("Removed client %s. %d registered clients", userID, len(s.clients))
		}
	}
}

func (s Sse) UseRouter(router *gin.RouterGroup) {
	router.StaticFile("/", "./infra/entrypoint/sse/index.html")
	router.GET(
		"/stream",
		s.middleware.Authenticate,
		s.registerClient,
		s.sendEvent,
	)
}

func (s Sse) ClientByUserID(userID string) Client {
	return s.clients[userID]
}

func (s Sse) registerClient(ctx *gin.Context) {
	userID := ctx.GetString(constant.UserIDContextKey)
	if userID == "" {
		return
	}

	client := make(Client)
	s.clients[userID] = client
	s.newClientUserIDs <- userID

	defer func() {
		s.closedClientUserIDs <- userID
	}()

	ctx.Set(constant.ClientContextKey, client)
	ctx.Next()
}

func (s Sse) sendEvent(ctx *gin.Context) {
	v, ok := ctx.Get(constant.ClientContextKey)
	if !ok {
		return
	}
	client, ok := v.(Client)
	if !ok {
		return
	}

	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-client; ok {
			ctx.SSEvent("message", msg)
			return true
		}
		return false
	})
}
