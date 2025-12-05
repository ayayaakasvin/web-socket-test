// Handlers that serves for main http server, accessed via handlers.Handler struct that contains necessary dependencies
package handlers

import (
	"web-socket-test/internal/models/core"
	"web-socket-test/internal/services/jwtservice"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	cache    core.Cache
	jwtM     *jwtservice.JWTService
	clientS  core.ClientService
	userRepo core.UserRepository
	chatS    core.ChatService
	chatRepo core.ChatHistoryStorage

	upg *websocket.Upgrader

	logger *logrus.Logger
}

func NewHTTPHandlers(logger *logrus.Logger,
	cache core.Cache,
	jwtM *jwtservice.JWTService,
	clm core.ClientService,
	chatS core.ChatService,
	chatRepo core.ChatHistoryStorage,
	ur core.UserRepository,
	upg *websocket.Upgrader,
) *Handlers {
	return &Handlers{
		logger: logger,

		cache:    cache,
		jwtM:     jwtM,
		clientS:  clm,
		userRepo: ur,
		chatS:    chatS,
		chatRepo: chatRepo,

		upg: upg,
	}
}
