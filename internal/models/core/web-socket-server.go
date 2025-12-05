package core

import "context"

type WebSocketServer interface {
	Start (ctx context.Context) error
	Close (ctx context.Context) error
}