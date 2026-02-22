package add

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
	ID:       "add",
	Title:    "Добавить",
	Help:     "",
	Internal: true,
	Factory:  New,
}

type state int

const (
	stateInitial state = iota
	stateChooseSeason
	stateChooseAction
	stateChooseTorrent
	stateWaitFile
)

const requestTimeout = 2 * time.Minute
const maxTorrents uint32 = 8

type addCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (d *addCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	listId, err := strconv.Atoi(ctx.Arguments[0])
	if err != nil {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
	ctx.Arguments = ctx.Arguments[1:]

	id := ctx.Arguments[0]
	ctx.Arguments = ctx.Arguments[1:]

	list := rms_library.List(listId)
	switch list {
	case rms_library.List_Favourites, rms_library.List_WatchList, rms_library.List_Archive:
		_, err := d.f.NewLists().Add(ctx, &rms_library.ListsAddRequest{List: list, Id: id}, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			d.l.Logf(logger.ErrorLevel, "Add failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}
		return true, command.ReplyText("Добавлено")
	default:
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &addCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "add"}),
	}
}
