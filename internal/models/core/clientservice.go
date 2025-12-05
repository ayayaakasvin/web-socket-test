package core

import (
	"web-socket-test/internal/models"
)

type ClientService interface {
	ClientActions

	Close() int
}

type ClientActions interface {
	Register(c *models.Client)
	Unregister(userId uint)
	Snapshot() []*models.Client
}
