package core

import (
	"context"
	"web-socket-test/internal/models"
)

type UserRepository interface {
	UserActions

	Close() error
}

type UserActions interface {
	GetPublicUserInfo(ctx context.Context, userID uint) (*models.User, error)
	GetPrivateUserInfo(ctx context.Context, userID uint) (*models.User, error)
}
