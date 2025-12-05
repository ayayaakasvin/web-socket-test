package core

import (
	"context"
	"web-socket-test/internal/models/dto"
)

type MessageWriter interface {
	SaveMessage(ctx context.Context, wm *dto.WBSMessage) error
}