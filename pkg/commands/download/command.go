package download

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
	ID:       "download",
	Title:    "Скачать",
	Help:     "",
	Internal: true,
	Factory:  New,
}

type state int

const (
	stateInitial state = iota
	stateChooseSeason
	stateChooseTorrent
	stateWaitFile
)

const requestTimeout = 2 * time.Minute
const maxTorrents uint32 = 8

type downloadCommand struct {
	f         servicemgr.ServiceFactory
	l         logger.Logger
	state     state
	stateMap  map[state]command.Handler
	download  command.Handler
	faster    bool
	watchlist bool
	id        string
	season    *uint
	torrents  []string
	mov       *rms_library.Movie
}

func (d *downloadCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	return d.stateMap[d.state](ctx)
}

func (d *downloadCommand) doInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
	switch ctx.Arguments[0] {
	case "auto":
		d.download = d.downloadAuto

	case "faster":
		d.download = d.downloadAuto
		d.faster = true

	case "select":
		d.download = d.downloadSelect

	case "file":
		d.download = d.downloadFile

	case "watchlist":
		d.download = d.downloadAuto
		d.watchlist = true

	default:
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	d.id = ctx.Arguments[1]

	result, err := d.f.NewMovies().Get(ctx, &rms_library.GetMovieRequest{ID: d.id})
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Retrieve info about media failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	if result.Result == nil {
		d.l.Log(logger.WarnLevel, "Movie not found")
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	mov := result.Result
	d.mov = result.Result

	if mov.Info.Type != rms_library.MovieType_TvSeries || mov.Info.Seasons == nil || ctx.Arguments[0] == "file" {
		return d.download(ctx)
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

func (d *downloadCommand) downloadAuto(ctx command.Context) (bool, []*communication.BotMessage) {
	req := &rms_library.DownloadMovieAutoRequest{
		Id:           d.id,
		Faster:       d.faster,
		UseWatchList: d.watchlist,
	}
	if d.season != nil {
		season := uint32(*d.season)
		req.Season = &season
	}

	resp, err := d.f.NewMovies().DownloadAuto(ctx, req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "request to library failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	if !resp.Found {
		return true, command.ReplyText("Не удалось найти подходящую раздачу")
	}

	if len(resp.Seasons) <= 1 {
		return true, command.ReplyText("Скачивание началось")
	}

	return true, command.ReplyText("Удалось найти сезоны " + formatSeasons(resp.Seasons) + ". Скачивание началось")
}

func (d *downloadCommand) downloadSelect(ctx command.Context) (bool, []*communication.BotMessage) {
	req := rms_library.FindMovieTorrentsRequest{
		Id:    d.id,
		Limit: maxTorrents,
	}
	if d.season != nil {
		season := uint32(*d.season)
		req.Season = &season
	}

	resp, err := d.f.NewMovies().FindTorrents(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	if len(resp.Results) == 0 {
		return true, command.ReplyText(command.NothingFound)
	}

	for _, t := range resp.Results {
		d.torrents = append(d.torrents, t.Id)
	}

	d.state = stateChooseTorrent
	return false, []*communication.BotMessage{formatTorrents(resp.Results)}
}

func (d *downloadCommand) downloadFile(ctx command.Context) (bool, []*communication.BotMessage) {
	d.state = stateWaitFile
	return false, command.ReplyText("Необходимо прислать торрент-файл с содержимым выбранного фильма/сериала")
}

func (d *downloadCommand) doChooseSeason(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return false, command.ReplyText("Необходимо выбрать сезон")
	}
	if ctx.Arguments[0] == "Все" {
		return d.download(ctx)
	}
	season, err := strconv.ParseUint(ctx.Arguments[0], 10, 8)
	if err != nil {
		return false, command.ReplyText("Неверно указан номер сезона")
	}
	s := uint(season)
	d.season = &s

	return d.download(ctx)
}

func (d *downloadCommand) doChooseTorrent(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return false, command.ReplyText("Необходимо выбрать раздачу")
	}
	no, err := strconv.ParseInt(ctx.Arguments[0], 10, 8)
	if err != nil || no <= 0 || no > int64(len(d.torrents)) {
		return false, command.ReplyText("Неверно указан номер раздачи")
	}

	id := d.torrents[no-1]

	_, err = d.f.NewMovies().Download(ctx, &rms_library.DownloadTorrentRequest{TorrentId: id}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Download request failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Скачивание началось")
}

func (d *downloadCommand) doWaitFile(ctx command.Context) (bool, []*communication.BotMessage) {
	if ctx.Attachment == nil {
		return false, command.ReplyText("Необходимо прислать торрент-файл")
	}
	if ctx.Attachment.MimeType != "application/x-bittorrent" {
		return false, command.ReplyText("Неверный формат файла")
	}

	req := rms_library.UploadMovieRequest{
		Id:          d.mov.Id,
		Info:        d.mov.Info,
		TorrentFile: ctx.Attachment.Content,
	}
	_, err := d.f.NewMovies().Upload(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Upload request failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Скачивание началось")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	d := &downloadCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "download"}),
	}

	d.stateMap = map[state]command.Handler{
		stateInitial:       d.doInitial,
		stateChooseSeason:  d.doChooseSeason,
		stateChooseTorrent: d.doChooseTorrent,
		stateWaitFile:      d.doWaitFile,
	}

	return d
}
