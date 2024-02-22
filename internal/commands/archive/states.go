package archive

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type state int

const (
	stateInitial state = iota
	stateChooseCamera
	stateChooseDay
	stateChooseTime
	stateChooseDuration
)

const requestTimeout = 30 * time.Second

func (c *archiveCommand) doInitial(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	c.state = stateChooseCamera
	if len(args) == 0 {
		list, err := c.f.NewCctv().GetCameras(ctx, &emptypb.Empty{}, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			c.l.Logf(logger.ErrorLevel, "Get cameras failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}

		for _, cam := range list.Cameras {
			c.cameras[cam.Name] = cam.Id
		}

		return false, []*communication.BotMessage{formatCameraList(list.Cameras)}
	}

	return c.fn[c.state](ctx, args, attachment)
}

func (c *archiveCommand) doChooseCamera(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	if len(args) == 0 {
		return false, command.ReplyText("Необходимо указать камеру")
	}

	ok := false
	c.camera, ok = c.cameras[args[0]]
	if !ok {
		return false, command.ReplyText("Неверно указана камера")
	}

	args = args[1:]
	c.state = stateChooseDay

	if len(args) == 0 {
		return false, []*communication.BotMessage{formatDayRequest()}
	}
	return c.fn[c.state](ctx, args, attachment)
}

func (c *archiveCommand) doChooseDay(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	if len(args) == 0 {
		return false, []*communication.BotMessage{formatDayRequest()}
	}

	t, ok := parseDay(args[0])
	if !ok {
		return false, command.ReplyText("Неверно введена дата")
	}

	c.ts = t
	args = args[1:]
	c.state = stateChooseTime

	if len(args) == 0 {
		return false, []*communication.BotMessage{formatTimeRequest()}
	}
	return c.fn[c.state](ctx, args, attachment)
}

func (c *archiveCommand) doChooseTime(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	if len(args) == 0 {
		return false, []*communication.BotMessage{formatTimeRequest()}
	}

	t, ok := parseTime(args[0])
	if !ok {
		return false, command.ReplyText("Неверно введено время")
	}

	c.ts.Add(t)
	args = args[1:]
	c.state = stateChooseDuration

	if len(args) == 0 {
		return false, []*communication.BotMessage{formatDurationRequest()}
	}
	return c.fn[c.state](ctx, args, attachment)
}

func (c *archiveCommand) doChooseDuration(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	if len(args) == 0 {
		return false, []*communication.BotMessage{formatDurationRequest()}
	}

	dur, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return false, command.ReplyText("Неверно введено время")
	}

	c.dur = dur
	if err = c.start(ctx); err != nil {
		c.l.Logf(logger.ErrorLevel, "Start archive download failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText("Выгрузка видео началась. По окончанию придет уведомление")
}
