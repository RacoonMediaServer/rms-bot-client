package downloads

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:      "downloads",
	Title:   "Загрузки",
	Help:    "Управление загрузками контента",
	Factory: New,
}

type downloadsCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (d *downloadsCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		return d.doList(ctx)
	}

	arg := ctx.Arguments[0]
	ctx.Arguments = ctx.Arguments[1:]
	switch arg {
	case "remove":
		return d.doRemove(ctx)
	case "up":
		return d.doUp(ctx)
	}

	return true, command.ReplyText(command.ParseArgumentsFailed)
}

func (d *downloadsCommand) doList(ctx command.Context) (bool, []*communication.BotMessage) {
	resp, err := d.f.NewTorrent().GetTorrents(ctx, &rms_torrent.GetTorrentsRequest{IncludeDoneTorrents: false})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Get torrents failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	if len(resp.Torrents) == 0 {
		return true, command.ReplyText("Нет активных загрузок")
	}
	var messages []*communication.BotMessage
	for _, t := range resp.Torrents {
		messages = append(messages, formatTorrent(t))
	}
	return true, messages
}

func (d *downloadsCommand) doRemove(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	_, err := d.f.NewTorrent().RemoveTorrent(ctx, &rms_torrent.RemoveTorrentRequest{Id: ctx.Arguments[0]})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Remove torrent failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText(command.Removed)
}

func (d *downloadsCommand) doUp(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return true, command.ReplyText("Не удалось распознать параметры команды")
	}

	_, err := d.f.NewTorrent().UpPriority(ctx, &rms_torrent.UpPriorityRequest{Id: ctx.Arguments[0]})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Up torrent failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText("Приоритет изменился")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &downloadsCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "downloads"}),
	}
}
