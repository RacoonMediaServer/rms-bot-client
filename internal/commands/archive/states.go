package archive

import (
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

func (c *archiveCommand) doInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	c.state = stateChooseCamera

	list, err := c.interlayer.Services.NewCctv().GetCameras(ctx.Ctx, &emptypb.Empty{}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Get cameras failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	for _, cam := range list.Cameras {
		c.cameras[cam.Name] = cam.Id
	}

	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatCameraList(list.Cameras)}
	}

	return c.fn[c.state](ctx)
}

func (c *archiveCommand) doChooseCamera(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		return false, command.ReplyText("Необходимо указать камеру")
	}

	ok := false
	c.camera, ok = c.cameras[ctx.Arguments[0]]
	if !ok {
		return false, command.ReplyText("Неверно указана камера")
	}
	c.ui.Camera = ctx.Arguments[0]

	ctx.Arguments = ctx.Arguments[1:]
	c.state = stateChooseDay

	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatDayRequest()}
	}
	return c.fn[c.state](ctx)
}

func (c *archiveCommand) doChooseDay(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatDayRequest()}
	}

	t, ok := parseDay(ctx.Arguments[0])
	if !ok {
		return false, command.ReplyText("Неверно введена дата")
	}

	c.ts = t
	ctx.Arguments = ctx.Arguments[1:]
	c.state = stateChooseTime

	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatTimeRequest()}
	}
	return c.fn[c.state](ctx)
}

func (c *archiveCommand) doChooseTime(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatTimeRequest()}
	}

	t, ok := parseTime(ctx.Arguments[0])
	if !ok {
		return false, command.ReplyText("Неверно введено время")
	}

	c.ts = c.ts.Add(t)
	c.ui.Time = c.ts.Local().Format(time.RFC3339)
	ctx.Arguments = ctx.Arguments[1:]
	c.state = stateChooseDuration

	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatDurationRequest()}
	}
	return c.fn[c.state](ctx)
}

func (c *archiveCommand) doChooseDuration(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		return false, []*communication.BotMessage{formatDurationRequest()}
	}

	dur, err := strconv.ParseUint(ctx.Arguments[0], 10, 16)
	if err != nil {
		return false, command.ReplyText("Неверно введено время")
	}

	c.dur = dur
	c.ui.Duration = uint(dur)
	if err = c.start(ctx.Ctx, ctx.UserID); err != nil {
		c.l.Logf(logger.ErrorLevel, "Start archive download failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText("Выгрузка видео началась. По окончанию придет уведомление")
}
