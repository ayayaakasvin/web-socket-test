package core

import (
	"context"

	"github.com/ayayaakasvin/web-socket-test/internal/models/dto"
)

type MessageWriter interface {
	SaveMessage(ctx context.Context, wm *dto.WBSMessage) error
}
