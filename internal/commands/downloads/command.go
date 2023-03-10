package downloads

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

import (
	"context"
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

func replyText(text string) []*communication.BotMessage {
	return []*communication.BotMessage{
		{
			Text: text,
		},
	}
}

func (d *downloadsCommand) Do(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) == 0 {
		return d.doList(ctx, arguments)
	}

	switch arguments[0] {
	case "remove":
		return d.doRemove(ctx, arguments[1:])
	case "up":
		return d.doUp(ctx, arguments[1:])
	}

	return true, replyText(command.ParseArgumentsFailed)
}

func (d *downloadsCommand) doList(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	resp, err := d.f.NewTorrent().GetTorrents(ctx, &rms_torrent.GetTorrentsRequest{IncludeDoneTorrents: false})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Get torrents failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}
	if len(resp.Torrents) == 0 {
		return true, replyText("Нет активных загрузок")
	}
	var messages []*communication.BotMessage
	for _, t := range resp.Torrents {
		messages = append(messages, formatTorrent(t))
	}
	return true, messages
}

func (d *downloadsCommand) doRemove(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) != 1 {
		return true, replyText(command.ParseArgumentsFailed)
	}

	_, err := d.f.NewTorrent().RemoveTorrent(ctx, &rms_torrent.RemoveTorrentRequest{Id: arguments[0]})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Remove torrent failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}
	return true, replyText(command.Removed)
}

func (d *downloadsCommand) doUp(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) != 1 {
		return true, replyText("Не удалось распознать параметры команды")
	}

	_, err := d.f.NewTorrent().UpPriority(ctx, &rms_torrent.UpPriorityRequest{Id: arguments[0]})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Up torrent failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}
	return true, replyText("Приоритет изменился")
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &downloadsCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "downloads"}),
	}
}
