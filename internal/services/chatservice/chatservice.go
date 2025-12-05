package chatservice

import (
	"context"
	"time"
	"web-socket-test/internal/models"
	"web-socket-test/internal/models/core"
	"web-socket-test/internal/models/dto"

	"github.com/sirupsen/logrus"
)

type ChatService struct {
	broadcast chan *dto.WBSMessage
	cm        core.ClientService
	logger    *logrus.Logger
	writer    core.MessageWriter
}

func NewChatService(cm core.ClientService, logger *logrus.Logger) *ChatService {
	return &ChatService{
		broadcast: make(chan *dto.WBSMessage, 256),
		cm:        cm,
		logger:    logger,
	}
}

// Context is passed through ShutdownService, to notify that main context is cancelled, time to terminate
func (cs *ChatService) RunBroadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-cs.broadcast:
			if !ok {
				cs.logger.Info("broadcast channel closed, stopping writer")
				return
			}

			clients := cs.cm.Snapshot()
			for _, cl := range clients {
				client := cl
				if err := cs.WriteMessageToClient(client, msg); err != nil {
					cs.logger.WithError(err).Warn("failed to write to client, closing the connection")
					client.Conn.Close()
					cs.cm.Unregister(client.UserInfo.ID)
					cs.PushMessage(dto.SystemMessage(client.UserInfo.ID, client.UserInfo.Username, dto.DisconnectType))
				}
			}

			cs.logger.WithField("msg", msg.Payload).WithField("username", msg.Origin.Username).Info("message was written to clients")

		}
	}
}

func (cs *ChatService) WriteMessageToClient(c *models.Client, msg *dto.WBSMessage) error {
	c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.Conn.WriteJSON(msg)
}

func (cs *ChatService) PushMessage(msg *dto.WBSMessage) {
    select {
    case cs.broadcast <- msg:
    default:
        cs.logger.WithField("len", len(cs.broadcast)).Warn("broadcast channel can not accept msg!")
    }
}