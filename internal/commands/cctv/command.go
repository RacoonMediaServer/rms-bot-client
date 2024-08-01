package cctv

import (
	"context"
	"fmt"
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const requestTimeout = 20 * time.Second

var Command command.Type = command.Type{
	ID:      "cctv",
	Title:   "Видеонаблюдение",
	Help:    "Управление системой видеонаблюдения",
	Factory: New,
}

type cctvCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (c cctvCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		msg := communication.BotMessage{}
		msg.Text = "Выберите команду из списка"
		msg.KeyboardStyle = communication.KeyboardStyle_Chat
		msg.Buttons = []*communication.Button{
			{
				Title:   nobodyAtHome,
				Command: nobodyAtHome,
			},
			{
				Title:   iamAtHome,
				Command: iamAtHome,
			},
		}
		return false, []*communication.BotMessage{&msg}
	}

	nobodyAtHomeMode := false

	cmd := ctx.Arguments.String()
	switch cmd {
	case nobodyAtHome:
		nobodyAtHomeMode = true
	case iamAtHome:
	default:
		return false, command.ReplyText("Не удалось распознать команду")
	}

	if err := c.setNobodyAtHomeMode(ctx, nobodyAtHomeMode); err != nil {
		c.l.Logf(logger.ErrorLevel, "Set 'nobodyAtHome' to %d failed: %s", nobodyAtHomeMode, err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	if nobodyAtHomeMode {
		msg := communication.BotMessage{
			Text: fmt.Sprintf("<b>%s</b>", nobodyAtHome),
			Pin:  communication.BotMessage_ThisMessage,
			User: command.BroadcastMessage,
		}
		return true, []*communication.BotMessage{&msg}
	}

	msg := communication.BotMessage{
		Text: fmt.Sprintf("<b>%s</b>", iamAtHome),
		Pin:  communication.BotMessage_Drop,
		User: command.BroadcastMessage,
	}
	return true, []*communication.BotMessage{&msg}
}

func (c cctvCommand) setNobodyAtHomeMode(ctx context.Context, active bool) error {
	_, err := c.f.NewCctvCameras().SetNobodyAtHomeMode(ctx, &rms_cctv.SetNobodyAtHomeModeRequest{Active: active}, client.WithRequestTimeout(requestTimeout))
	return err
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	nc := &cctvCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "cctv"}),
	}

	return nc
}
