package remove

import (
	"strconv"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "remove",
	Title:    "Переместить",
	Help:     "Перемещение между списка",
	Factory:  New,
	Internal: true,
}

type removeCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (r *removeCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
	id := ctx.Arguments[0]
	listInt, err := strconv.ParseInt(ctx.Arguments[1], 10, 32)
	if err != nil {
		r.l.Logf(logger.ErrorLevel, "Parse list failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	list := rms_library.List(listInt)

	if _, err := r.f.NewLists().Move(ctx, &rms_library.ListsMoveRequest{Id: id, List: list}); err != nil {
		r.l.Logf(logger.ErrorLevel, "Move movie failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Перемещено")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &removeCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "remove"}),
	}
}
