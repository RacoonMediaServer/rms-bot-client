package download

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"strconv"
	"time"
)

var Command command.Type = command.Type{
	ID:       "download",
	Title:    "Скачать",
	Help:     "",
	Internal: true,
	Factory:  New,
}

type doFunc func(ctx context.Context, args command.Arguments) (bool, []*communication.BotMessage)
type state int

const (
	stateInitial state = iota
	stateChooseSeason
	stateChooseTorrent
)

const requestTimeout = 1 * time.Minute
const maxTorrents uint32 = 8

type downloadCommand struct {
	f        servicemgr.ServiceFactory
	l        logger.Logger
	state    state
	stateMap map[state]doFunc
	download doFunc
	id       string
	season   *uint
	torrents []string
}

func replyText(text string) []*communication.BotMessage {
	return []*communication.BotMessage{
		{
			Text: text,
		},
	}
}

func (d *downloadCommand) Do(ctx context.Context, arguments command.Arguments) (done bool, messages []*communication.BotMessage) {
	return d.stateMap[d.state](ctx, arguments)
}

func (d *downloadCommand) doInitial(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) < 2 {
		return true, replyText(command.ParseArgumentsFailed)
	}
	switch arguments[0] {
	case "auto":
		d.download = d.downloadAuto
	case "select":
		d.download = d.downloadSelect
	default:
		return true, replyText(command.ParseArgumentsFailed)
	}

	d.id = arguments[1]

	result, err := d.f.NewLibrary().GetMovie(ctx, &rms_library.GetMovieRequest{ID: d.id})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Retrieve info about media failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}
	if result.Result == nil {
		d.l.Log(logger.WarnLevel, "Movie not found")
		return true, replyText(command.SomethingWentWrong)
	}

	mov := result.Result

	if mov.Info.Type == rms_library.MovieType_Film || mov.Info.Seasons == nil {
		return d.download(ctx, arguments[1:])
	}

	d.state = stateChooseSeason
	msg := communication.BotMessage{Text: "Выберите сезон"}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Все",
		Command: "Все",
	})

	for i := uint32(1); i <= *mov.Info.Seasons; i++ {
		if mov.TvSeries != nil {
			_, ok := mov.TvSeries.Seasons[i]
			if ok {
				continue
			}
		}
		no := strconv.FormatUint(uint64(i), 10)
		msg.Buttons = append(msg.Buttons, &communication.Button{Title: no, Command: no})
	}

	return false, []*communication.BotMessage{&msg}
}

func (d *downloadCommand) downloadAuto(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	req := &rms_library.DownloadMovieAutoRequest{Id: d.id}
	if d.season != nil {
		season := uint32(*d.season)
		req.Season = &season
	}

	resp, err := d.f.NewLibrary().DownloadMovieAuto(ctx, req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "request to library failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}

	if !resp.Found {
		return true, replyText("Не удалось найти подходящую раздачу")
	}

	if len(resp.Seasons) <= 1 {
		return true, replyText("Скачивание началось")
	}

	return true, replyText("Удалось найти сезоны " + formatSeasons(resp.Seasons) + ". Скачивание началось")
}

func (d *downloadCommand) downloadSelect(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	req := rms_library.FindMovieTorrentsRequest{
		Id:    d.id,
		Limit: maxTorrents,
	}
	if d.season != nil {
		season := uint32(*d.season)
		req.Season = &season
	}

	resp, err := d.f.NewLibrary().FindMovieTorrents(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return true, replyText(command.SomethingWentWrong)
	}
	if len(resp.Results) == 0 {
		return true, replyText(command.NothingFound)
	}

	for _, t := range resp.Results {
		d.torrents = append(d.torrents, t.Id)
	}

	d.state = stateChooseTorrent
	return false, []*communication.BotMessage{formatTorrents(resp.Results)}
}

func (d *downloadCommand) doChooseSeason(ctx context.Context, args command.Arguments) (bool, []*communication.BotMessage) {
	if len(args) != 1 {
		return false, replyText("Необходимо выбрать сезон")
	}
	if args[0] == "Все" {
		return d.download(ctx, args)
	}
	season, err := strconv.ParseUint(args[0], 10, 8)
	if err != nil {
		return false, replyText("Неверно указан номер сезона")
	}
	s := uint(season)
	d.season = &s

	return d.download(ctx, args)
}

func (d *downloadCommand) doChooseTorrent(ctx context.Context, args command.Arguments) (bool, []*communication.BotMessage) {
	if len(args) != 1 {
		return false, replyText("Необходимо выбрать раздачу")
	}
	no, err := strconv.ParseInt(args[0], 10, 8)
	if err != nil || no <= 0 || no > int64(len(d.torrents)) {
		return false, replyText("Неверно указан номер раздачи")
	}

	id := d.torrents[no-1]

	_, err = d.f.NewLibrary().DownloadTorrent(ctx, &rms_library.DownloadTorrentRequest{TorrentId: id}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Download request failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}

	return true, replyText("Скачивание началось")
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	d := &downloadCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "download"}),
	}

	d.stateMap = map[state]doFunc{
		stateInitial:       d.doInitial,
		stateChooseSeason:  d.doChooseSeason,
		stateChooseTorrent: d.doChooseTorrent,
	}

	return d
}
