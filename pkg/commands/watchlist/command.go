package watchlist

import (
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const requestTimeout = 2 * time.Minute

var Command command.Type = command.Type{
	ID:       "watchlist",
	Title:    "Отложенное",
	Help:     "Управление отложенным медиа-контентом",
	Internal: true,
	Factory:  New,
}

type watchListCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (c watchListCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 || ctx.Arguments[0] != "add" {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	id := ctx.Arguments[1]
	_, err := c.f.NewMovies().WatchLater(ctx, &rms_library.WatchLaterRequest{Id: id}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Add to watchlist failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Добавлено. В списке отложенных контент появится с опозданием")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	nc := &watchListCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "watchlist"}),
	}

	return nc
}
