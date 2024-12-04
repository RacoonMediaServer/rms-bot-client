package main

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-bot-client/internal/cmdset"
	"github.com/RacoonMediaServer/rms-bot-client/internal/config"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/background"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/bot"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/session"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/speech"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	// Plugins
	_ "github.com/go-micro/plugins/v4/registry/etcd"
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

	backService := background.NewService()
	defer backService.Stop()

	interlayer := command.Interlayer{
		Services:       servicemgr.NewServiceFactory(service),
		TaskService:    backService,
		ShareDirectory: cfg.ContentDirectory,
	}

	serverSession := session.New(cfg.Remote, cfg.Device)
	defer serverSession.Shutdown()

	botSettings := bot.Settings{
		Transport:  serverSession,
		Interlayer: interlayer,
		CmdFactory: cmdset.New(),
	}
	if cfg.VoiceRecognition {
		botSettings.SpeechRecognizer = &speech.Recognizer{ServiceFactory: interlayer.Services}
	}

	botInstance := bot.New(botSettings)
	defer botInstance.Shutdown()

	// регистрируем хендлеры
	if err := rms_bot_client.RegisterRmsBotClientHandler(service.Server(), botInstance); err != nil {
		logger.Fatalf("Register service failed: %s", err)
	}

	if err := service.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}
}
