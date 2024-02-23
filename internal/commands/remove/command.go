package remove

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "remove",
	Title:    "Удалить",
	Help:     "Удаление фильмов и сериалов",
	Factory:  New,
	Internal: true,
}

type removeCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (r *removeCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	if _, err := r.f.NewLibrary().DeleteMovie(ctx, &rms_library.DeleteMovieRequest{ID: ctx.Arguments[0]}); err != nil {
		r.l.Logf(logger.ErrorLevel, "Remove movie failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText(command.Removed)
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &removeCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "remove"}),
	}
}
