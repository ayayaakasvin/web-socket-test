package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"web-socket-test/internal/config"
	"web-socket-test/internal/logger"
	"web-socket-test/internal/repository/postgresql"
	"web-socket-test/internal/repository/valkey"
	"web-socket-test/internal/services/chatservice"
	"web-socket-test/internal/services/clientservice"
	"web-socket-test/internal/services/jwtservice"
	websocketserver "web-socket-test/internal/web-socket-server"

	"github.com/ayayaakasvin/go-shutdown-channel"
	"github.com/sirupsen/logrus"
)

func main() {
    mainCtx, cancel := context.WithCancel(context.Background())
    wg := &sync.WaitGroup{}
    wg.Add(1)

    s:= goshutdownchannel.NewShutdown(mainCtx, cancel)
    s.Notify(os.Interrupt, syscall.SIGTERM)

    // Config and Logger
    cfg := config.MustLoadConfig()
    logger := logger.SetupLogger("WS")

    // Dependencies
    cc := valkey.NewValkeyClient(cfg.ValkeyConfig, s)
    chatRepo := postgresql.NewPostgreSQL_Mock()
    jwtS := jwtservice.NewJWTService(&cfg.JWTSecret)
    clientS := clientservice.NewClientService(logger)
    chatS := chatservice.NewChatService(clientS, logger)

    wbs := websocketserver.NewWBS(cfg, logger, chatRepo, clientS, chatS, jwtS, cc, chatRepo)

    go func () {
        defer wg.Done()

        <- s.Done()
		fmt.Println(s.Message())

        gracefulCtx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
        defer cancel()

        wbs.Close(gracefulCtx)
        logger.WithError(mainCtx.Err()).Info("Gracefully Shutdowning...")
    }()

    go func () {
        if err := wbs.Start(s.Context()); err != nil {
            logger.WithField("wbs error", err).Error("WBS Start")
            s.Send("wbs-server", err.Error())
            return
        }
    }()

    go PrintServerStatus(mainCtx, logger)

    wg.Wait()
    logger.WithError(mainCtx.Err()).Info("Graceful Shutdown completed")
}

func PrintServerStatus(ctx context.Context, logger *logrus.Logger) {
    ticker := time.NewTicker(time.Second * 60)

    for {
        select {
        case <- ticker.C:
            logger.Info("Server is alive...")
        case <- ctx.Done():
            return
        }
    }
}