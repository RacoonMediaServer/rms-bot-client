package snapshot

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const requestTimeout = 20 * time.Second

var Command = command.Type{
	ID:      "snapshot",
	Title:   "Снепшот",
	Help:    "Получить снепшот с камеры",
	Factory: New,
}

type snapshotCommand struct {
	f           servicemgr.ServiceFactory
	l           logger.Logger
	mapNameToId map[string]uint32
}

func (c *snapshotCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	switch len(ctx.Arguments) {
	case 0:
		return c.doListCameras(ctx)
	case 1:
		return c.doSnapshot(ctx, ctx.Arguments[0])
	default:
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
}

func (c *snapshotCommand) doListCameras(ctx context.Context) (bool, []*communication.BotMessage) {
	list, err := c.f.NewCctv().GetCameras(ctx, &emptypb.Empty{}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Get cameras failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	for _, cam := range list.Cameras {
		c.mapNameToId[cam.Name] = cam.Id
	}

	return false, []*communication.BotMessage{formatCameraList(list.Cameras)}
}

func (c *snapshotCommand) doSnapshot(ctx context.Context, cameraName string) (bool, []*communication.BotMessage) {
	id, ok := c.mapNameToId[cameraName]
	if !ok {
		return true, command.ReplyText("Камера с таким именем не найдена")
	}

	resp, err := c.f.NewCctv().GetSnapshot(ctx, &rms_cctv.GetSnapshotRequest{CameraId: id}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Get snapshot failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, []*communication.BotMessage{formatSnapshot(cameraName, resp.Snapshot)}
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &snapshotCommand{
		f:           f,
		l:           l.Fields(map[string]interface{}{"command": "snapshot"}),
		mapNameToId: make(map[string]uint32),
	}
}
