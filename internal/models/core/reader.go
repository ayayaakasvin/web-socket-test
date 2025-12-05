package core

import (
	"context"
	"web-socket-test/internal/models/dto"
)

type MessageReader interface {
	GetRecentMessage(ctx context.Context, limit int) ([]*dto.WBSMessage, error)
}