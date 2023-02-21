package main

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-bot-client/internal/bot"
	"github.com/RacoonMediaServer/rms-bot-client/internal/config"
	"github.com/RacoonMediaServer/rms-bot-client/internal/session"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var Version = "v0.0.0"

const serviceName = "rms-bot-client"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	cfg := config.Config()
	serverSession := session.New(cfg.Remote, cfg.Device)
	defer serverSession.Shutdown()

	botInstance := bot.New(serverSession, servicemgr.NewServiceFactory(service))
	defer botInstance.Shutdown()

	// регистрируем хендлеры
	if err := rms_bot_client.RegisterRmsBotClientHandler(service.Server(), botInstance); err != nil {
		logger.Fatalf("Register service failed: %s", err)
	}

	if err := service.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}
}
