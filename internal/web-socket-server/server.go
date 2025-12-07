package websocketserver

import (
	"context"
	"net/http"
	"sync"

	"github.com/ayayaakasvin/web-socket-test/internal/config"
	"github.com/ayayaakasvin/web-socket-test/internal/models"
	"github.com/ayayaakasvin/web-socket-test/internal/models/core"
	"github.com/ayayaakasvin/web-socket-test/internal/models/dto"
	"github.com/ayayaakasvin/web-socket-test/internal/services/jwtservice"
	"github.com/ayayaakasvin/web-socket-test/internal/web-socket-server/handlers"
	"github.com/ayayaakasvin/web-socket-test/internal/web-socket-server/middlewares"

	"github.com/ayayaakasvin/lightmux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WBS struct {
	server *http.Server
	lmux   *lightmux.LightMux

	chatRepo core.ChatHistoryStorage
	clientS  core.ClientService
	chatS    core.ChatService
	jwtS     *jwtservice.JWTService
	cc       core.Cache
	ur       core.UserRepository
	corsCfg  *config.CorsConfig

	logger *logrus.Logger

	wbsUpg     *websocket.Upgrader
	wbsClients map[models.ID]*websocket.Conn

	mutex     *sync.RWMutex
	broadcast chan *dto.WBSMessage
}

func NewWBS(cfg *config.Config, logger *logrus.Logger,
	chatRepo core.ChatHistoryStorage,
	clientS core.ClientService,
	chatS core.ChatService,
	jwtManager *jwtservice.JWTService,
	cc core.Cache,
	ur core.UserRepository,
) core.WebSocketServer {
	wbs := &WBS{
		server: &http.Server{
			Addr:         cfg.Address,
			IdleTimeout:  cfg.IdleTimeout,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
		wbsUpg: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// the CheckOrigin is set up inside WBS.setupLightMux method 
		},

		jwtS:     jwtservice.NewJWTService(&cfg.JWTSecret),
		chatRepo: chatRepo,
		clientS:  clientS,
		chatS:    chatS,
		cc:       cc,
		ur:       ur,
		logger:   logger,

		corsCfg: &cfg.CorsConfig,

		wbsClients: make(map[models.ID]*websocket.Conn),
		mutex:      new(sync.RWMutex),
		broadcast:  make(chan *dto.WBSMessage),
	}
	wbs.lmux = lightmux.NewLightMux(wbs.server)

	return wbs
}

func (wbs *WBS) setupLightMux() {
	hndlrs := handlers.NewHTTPHandlers(wbs.logger, wbs.cc, wbs.jwtS, wbs.clientS, wbs.chatS, wbs.chatRepo, wbs.ur, wbs.wbsUpg)
	mds := middlewares.NewHTTPMiddlewares(wbs.logger, wbs.corsCfg, wbs.cc, wbs.ur, wbs.jwtS)

	wbs.wbsUpg.CheckOrigin = mds.WebSocketCheckOrigin()

	wbs.lmux.Use(mds.RecoverMiddleware, mds.LoggerMiddleware, mds.CORSMiddleware)

	wbs.lmux.NewRoute("/ping").Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	wbs.lmux.NewRoute("/ws").Handle(http.MethodGet, hndlrs.WS_Handler())
	wbs.lmux.NewRoute("/ping").Handle(http.MethodGet, hndlrs.PingHandler())
	wbs.lmux.NewRoute("/chat-history", mds.JWTAuthMiddleware).Handle(http.MethodGet, hndlrs.GetChatHistory())
	wbs.lmux.NewRoute("/clients", mds.JWTAuthMiddleware, mds.JWTAdminMiddleware).Handle(http.MethodGet, hndlrs.GetClientList())

	wbs.logger.Info("LightMux has been set up")
}

func (wbs *WBS) Start(ctx context.Context) error {
	wbs.setupLightMux()
	go wbs.chatS.RunBroadcast(ctx)

	wbs.logger.Infof("Server has been started on port: %s", wbs.server.Addr)
	wbs.logger.Infof("Available handlers:\n")

	wbs.lmux.PrintMiddlewareInfo()
	wbs.lmux.PrintRoutes()

	wbs.logger.Info("Server has been set up")
	return wbs.lmux.RunContext(ctx)
}

func (wbs *WBS) Close(ctx context.Context) error {
	wbs.logger.Info("Closing down")

	wbs.clientS.Close()
	wbs.chatRepo.Close()
	wbs.cc.Close()

	return wbs.server.Shutdown(ctx)
}

// func (wbs *WBS) pingHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("pong"))
// 	}
// }

// func (wbs *WBS) wsHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		token := r.URL.Query().Get("token")
// 		if token == "" {
// 			response.SendErrorJson(w, http.StatusUnauthorized, "missing token")
// 			wbs.logger.Warn("WebSocket connection attempt without token")
// 			return
// 		}

// 		claims, err := wbs.jwtS.ValidateJWT(token)
// 		if err != nil {
// 			response.SendErrorJson(w, http.StatusUnauthorized, "invalid token")
// 			wbs.logger.WithError(err).Warn("WebSocket token validation failed")
// 			return
// 		}

// 		userIdAny, ok := claims["user_id"]
// 		if !ok {
// 			response.SendErrorJson(w, http.StatusUnauthorized, "user_id missing")
// 			return
// 		}

// 		userId, err := wbs.jwtS.FetchUserID(userIdAny)
// 		if err != nil {
// 			response.SendErrorJson(w, http.StatusUnauthorized, "invalid user_id")
// 			return
// 		}

// 		conn, err := wbs.wbsUpg.Upgrade(w, r, nil)
// 		if err != nil {
// 			wbs.logger.WithError(err).Error("WebSocket upgrade failed")
// 			return
// 		}

// 		id := models.NewID(r)
// 		wbs.logger.WithField("addr", conn.RemoteAddr().String()).Info("Client connected")

// 		wbs.mutex.Lock()
// 		wbs.registerConnection(conn, id)
// 		wbs.mutex.Unlock()

// 		go wbs.readCLientMessage(conn, id, userId)
// 	}
// }

// func (wbs *WBS) readCLientMessage(conn *websocket.Conn, id models.ID, userId uint) {
// 	defer func() {
// 		wbs.mutex.Lock()
// 		conn.Close()
// 		delete(wbs.wbsClients, id)
// 		wbs.mutex.Unlock()
// 		wbs.logger.WithField("id", id).Info("Connection closed")
// 	}()

// 	for {
// 		msg := new(dto.WBSMessage)
// 		if err := conn.ReadJSON(msg); err != nil {
// 			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
// 				wbs.logger.WithField("id", id).Info("Client disconnected")
// 			} else {
// 				wbs.logger.WithError(err).Warn("Read error")
// 			}
// 			return
// 		}

// 		msg.Origin.UserID = int(userId)
// 		msg.Time = time.Now().Format("15:04")

// 		select {
// 		case wbs.broadcast <- msg:
// 		default:
// 			wbs.logger.Warn("Broadcast channel full, message dropped")
// 		}
// 	}
// }

// func (wbs *WBS) writeToClients() {
// 	for {
// 		select {
// 		case msg, ok := <-wbs.broadcast:
// 			if !ok {
// 				log.Info("broadcast channel closed, stopping writer")
// 				return
// 			}

// 			wbs.mutex.Lock()
// 			for id, conn := range wbs.wbsClients {
// 				c := conn
// 				c.SetWriteDeadline(time.Now().Add(1 * time.Second))
// 				if err := c.WriteJSON(msg); err != nil {
// 					log.WithError(err).Warn("failed to write to client, removing")
// 					c.Close()
// 					delete(wbs.wbsClients, id)
// 				}
// 			}
// 			wbs.mutex.Unlock()
// 			log.WithField("msg", msg.Payload).Info("message was written to clients")

// 		case <-time.After(time.Second * 1):
// 			continue
// 		}
// 	}
// }

// func (wbs *WBS) registerConnection(conn *websocket.Conn, id models.ID) {
// 	if oldConn, exists := wbs.wbsClients[id]; exists {
// 		oldConn.Close()
// 	}

// 	wbs.wbsClients[id] = conn
// }
