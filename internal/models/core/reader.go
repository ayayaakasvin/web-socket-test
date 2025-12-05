package core

import (
	"context"

	"github.com/ayayaakasvin/web-socket-test/internal/models/dto"
)

type MessageReader interface {
	GetRecentMessage(ctx context.Context, limit int) ([]*dto.WBSMessage, error)
}
