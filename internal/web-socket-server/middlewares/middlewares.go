package middlewares

import (
	"strings"

	"github.com/ayayaakasvin/web-socket-test/internal/config"
	"github.com/ayayaakasvin/web-socket-test/internal/models/core"
	"github.com/ayayaakasvin/web-socket-test/internal/services/jwtservice"

	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	logger     *logrus.Logger
	cache      core.Cache
	userRepo   core.UserRepository
	jwtManager *jwtservice.JWTService

	allowedOrigins   string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool

	allowedOriginsMap map[string]bool
}

func NewHTTPMiddlewares(logger *logrus.Logger, corsCfg *config.CorsConfig, cache core.Cache, userRepo core.UserRepository, jwtManager *jwtservice.JWTService) *Middlewares {
	mp := make(map[string]bool)
	for _, origin := range corsCfg.AllowedOrigins {
		mp[origin] = true
	}
	
	return &Middlewares{
		logger:     logger,
		cache:      cache,
		userRepo:   userRepo,
		jwtManager: jwtManager,

		allowedOrigins:   strings.Join(corsCfg.AllowedOrigins, ","),
		allowedMethods:   strings.Join(corsCfg.AllowedMethods, ","),
		allowedHeaders:   strings.Join(corsCfg.AllowedHeaders, ","),
		allowCredentials: corsCfg.AllowedCredentials,

		allowedOriginsMap: make(map[string]bool),
	}
}
