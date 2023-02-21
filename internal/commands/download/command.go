package download

import (
	"context"
	"fmt"
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

const requestTimeout = 1 * time.Minute

type downloadCommand struct {
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

func formatSeasons(seasons []uint32) string {
	result := ""
	for _, s := range seasons {
		result += fmt.Sprintf("%d, ", s)
	}
	if len(result) > 2 {
		result = result[0 : len(result)-2]
	}
	return result
}

func (d downloadCommand) Do(ctx context.Context, arguments command.Arguments) (done bool, messages []*communication.BotMessage) {
	if len(arguments) < 1 {
		return true, replyText("Параметры команды не распознаны")
	}
	req := &rms_library.DownloadMovieRequest{Id: arguments[0]}
	if len(arguments) >= 2 {
		// второй параметр парси как номер сезона
		season, err := strconv.ParseUint(arguments[1], 10, 8)
		if err != nil {
			return true, replyText("Неверно указан номер сезона")
		}
		no := uint32(season)
		req.Season = &no
	}

	resp, err := d.f.NewLibrary().DownloadMovie(ctx, req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "request to library failed: %s", err)
		return true, replyText("Что-то пошло не так...")
	}

	if !resp.Found {
		return true, replyText("Не удалось найти подходящую раздачу")
	}

	if len(resp.Seasons) == 0 {
		return true, replyText("Скачивание началось")
	}

	return true, replyText("Удалось найти сезоны " + formatSeasons(resp.Seasons) + ". Скачивание началось")
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &downloadCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "download"}),
	}
}
