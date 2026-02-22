package torrents

import (
	"strconv"
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "torrents",
	Title:    "Раздачи",
	Help:     "Управление раздачами",
	Factory:  New,
	Internal: true,
}

type torrentsCommand struct {
	f        servicemgr.ServiceFactory
	l        logger.Logger
	state    state
	id       string
	torrents []*rms_library.Torrent
}

type state int

const (
	stateInitial state = iota
	stateChoseSeason
	stateChooseTorrent
	stateWaitFile
)

func (t *torrentsCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch t.state {
	case stateInitial:
		return t.doStateInitial(ctx)
	case stateChoseSeason:
		return t.doSeasonChoice(ctx)
	case stateChooseTorrent:
		return t.doTorrentChoice(ctx)
	case stateWaitFile:
		return t.doWaitFile(ctx)
	default:
		return true, command.ReplyText(command.SomethingWentWrong)
	}
}

func (t *torrentsCommand) doStateInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 1 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	id := ctx.Arguments[0]
	ctx.Arguments = ctx.Arguments[1:]

	if len(ctx.Arguments) == 0 {
		return t.listTorrents(ctx, id)
	}

	cmd := ctx.Arguments[0]
	ctx.Arguments = ctx.Arguments[1:]

	switch cmd {
	case "remove":
		return t.removeTorrent(ctx, id, ctx.Arguments.String())
	case "add":
		return t.addTorrent(ctx, id)
	case "file":
		t.id = id
		t.state = stateWaitFile
		return false, command.ReplyText("Пришлите файл раздачи")
	default:
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
}

func (t *torrentsCommand) listTorrents(ctx command.Context, id string) (bool, []*communication.BotMessage) {
	torrentsService := t.f.NewTorrents()

	resp, err := torrentsService.List(ctx, &rms_library.TorrentsListRequest{Id: id})
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Torrents list failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	messages := []*communication.BotMessage{newAddMessage(id)}

	for _, torrent := range resp.Torrents {
		messages = append(messages, newTorrentMessage(id, torrent))
	}

	return true, messages
}

func (t *torrentsCommand) removeTorrent(ctx command.Context, id string, torrentId string) (bool, []*communication.BotMessage) {
	torrentsService := t.f.NewTorrents()

	_, err := torrentsService.Delete(ctx, &rms_library.TorrentsDeleteRequest{Id: id, TorrentId: torrentId})
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Torrents remove failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Раздача удалена")
}

func (t *torrentsCommand) addTorrent(ctx command.Context, id string) (bool, []*communication.BotMessage) {
	moviesService := t.f.NewMovies()
	skipSeasonsChoice := false

	resp, err := moviesService.Get(ctx, &rms_library.MoviesGetRequest{Id: id})
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Movies get failed: %s", err)
		skipSeasonsChoice = true
	} else {
		skipSeasonsChoice = resp.Info.Type != rms_library.MovieType_TvSeries || resp.Info.Seasons == nil
	}

	if !skipSeasonsChoice {
		t.state = stateChoseSeason
		t.id = id
		msg := communication.BotMessage{Text: "Выберите сезон"}
		msg.KeyboardStyle = communication.KeyboardStyle_Chat
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "Все",
			Command: "Все",
		})

		for i := uint32(1); i <= *resp.Info.Seasons; i++ {
			no := strconv.FormatUint(uint64(i), 10)
			msg.Buttons = append(msg.Buttons, &communication.Button{Title: no, Command: no})
		}

		return false, []*communication.BotMessage{&msg}
	}

	return t.findTorrents(ctx, id, nil)
}

func (t *torrentsCommand) findTorrents(ctx command.Context, id string, season *uint32) (bool, []*communication.BotMessage) {
	torrentsService := t.f.NewTorrents()
	resp, err := torrentsService.Find(ctx, &rms_library.TorrentsFindRequest{Id: id, Season: season}, client.WithRequestTimeout(1*time.Minute))
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Torrents find failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	if len(resp.Torrents) == 0 {
		return true, command.ReplyText("Раздачи не найдены")
	}

	t.state = stateChooseTorrent
	t.id = id
	t.torrents = resp.Torrents
	return false, []*communication.BotMessage{formatTorrents(resp.Torrents)}
}

func (t *torrentsCommand) doSeasonChoice(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 1 {
		return false, command.ReplyText(command.ParseArgumentsFailed)
	}

	seasonArg := ctx.Arguments[0]

	var season *uint32
	if seasonArg != "Все" {
		seasonInt, err := strconv.ParseUint(seasonArg, 10, 32)
		if err != nil {
			t.l.Logf(logger.ErrorLevel, "Parse season failed: %s", err)
			return false, command.ReplyText(command.ParseArgumentsFailed)
		}
		seasonUint := uint32(seasonInt)
		season = &seasonUint
	}

	return t.findTorrents(ctx, t.id, season)
}

func (t *torrentsCommand) doTorrentChoice(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 1 {
		return false, command.ReplyText(command.ParseArgumentsFailed)
	}
	torrentArg := ctx.Arguments[0]
	torrentInt, err := strconv.ParseInt(torrentArg, 10, 32)
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Parse torrent failed: %s", err)
		return false, command.ReplyText(command.ParseArgumentsFailed)
	}

	if torrentInt < 0 || torrentInt >= int64(len(t.torrents)) {
		return false, command.ReplyText(command.ParseArgumentsFailed)
	}

	torrent := t.torrents[torrentInt]
	_, err = t.f.NewTorrents().Add(ctx, &rms_library.TorrentsAddRequest{Id: t.id, Link: &torrent.Id})
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Torrents add failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Раздача добавлена")
}

func (t *torrentsCommand) doWaitFile(ctx command.Context) (bool, []*communication.BotMessage) {
	if ctx.Attachment == nil {
		return false, command.ReplyText("Пришлите файл раздачи")
	}

	if ctx.Attachment.MimeType != "application/x-bittorrent" {
		return false, command.ReplyText("Неверный тип файла, ожидается .torrent")
	}

	content := ctx.Attachment.Content

	_, err := t.f.NewTorrents().Add(ctx, &rms_library.TorrentsAddRequest{Id: t.id, TorrentFile: content})
	if err != nil {
		t.l.Logf(logger.ErrorLevel, "Torrents add by file failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Раздача добавлена")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &torrentsCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "remove"}),
	}
}
