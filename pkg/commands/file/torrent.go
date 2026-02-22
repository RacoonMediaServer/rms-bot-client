package file

import (
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const requestTimeout = 30 * time.Second

type tfhState int

const (
	tfhStateInit = iota
	tfhStateSelectTitle
)

type torrentFileHandler struct {
	f       servicemgr.ServiceFactory
	l       logger.Logger
	state   tfhState
	content []byte
}

func (t *torrentFileHandler) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch t.state {
	case tfhStateInit:
		t.content = ctx.Attachment.Content
		t.state = tfhStateSelectTitle
		if len(ctx.Arguments) == 0 {
			return false, command.ReplyText("Введите название фильма/сериала")
		}
		fallthrough

	case tfhStateSelectTitle:
		if len(ctx.Arguments) == 0 {
			return false, command.ReplyText("Название не должно быть пустым")
		}

		libraryService := t.f.NewMovies()

		req := rms_library.MoviesAddClipRequest{
			Title:   &ctx.Arguments[0],
			Torrent: t.content,
			List:    rms_library.List_WatchList,
		}

		_, err := libraryService.AddClip(ctx, &req, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			t.l.Logf(logger.ErrorLevel, "Upload clip failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}

		return true, command.ReplyText("Добавлено к просмотру")
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}
