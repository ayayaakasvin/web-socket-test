package middlewares

import (
	"github.com/ayayaakasvin/web-socket-test/internal/models/core"
	"github.com/ayayaakasvin/web-socket-test/internal/services/jwtservice"

	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	logger     *logrus.Logger
	cache      core.Cache
	userRepo   core.UserRepository
	jwtManager *jwtservice.JWTService
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache core.Cache, userRepo core.UserRepository, jwtManager *jwtservice.JWTService) *Middlewares {
	return &Middlewares{
		logger:     logger,
		cache:      cache,
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}
