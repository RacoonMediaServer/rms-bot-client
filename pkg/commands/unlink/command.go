package unlink

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:      "unlink",
	Title:   "Отключиться",
	Help:    "Отключить текущего пользователя Telegram от управления устройством",
	Factory: New,
}

type unlinkCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (u *unlinkCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	messages := []*communication.BotMessage{
		{
			Text: "Управление отключено",
		},
		{
			Type: communication.MessageType_UnlinkUser,
		},
	}

	return true, messages
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &unlinkCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "unlink"}),
	}
}
