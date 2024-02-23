package file

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/google/uuid"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"time"
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

		libraryService := t.f.NewLibrary()
		req := rms_library.UploadMovieRequest{
			Id: "internal:" + uuid.NewString(),
			Info: &rms_library.MovieInfo{
				Title: ctx.Arguments.String(),
				Type:  rms_library.MovieType_Clip,
			},
			TorrentFile: t.content,
		}

		_, err := libraryService.UploadMovie(ctx, &req, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			t.l.Logf(logger.ErrorLevel, "Upload movie failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}

		return true, command.ReplyText("Загрузка началась")
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}
