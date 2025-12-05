package core

import (
	"context"

	"github.com/ayayaakasvin/web-socket-test/internal/models/dto"
)

type ChatService interface {
	ChatActions
}

type ChatActions interface {
	RunBroadcast(ctx context.Context)
	PushMessage(msg *dto.WBSMessage)
}
