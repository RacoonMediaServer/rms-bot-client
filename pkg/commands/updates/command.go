package updates

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var Command command.Type = command.Type{
	ID:      "updates",
	Title:   "Обновления",
	Help:    "Запрос информации о новых сезонах сериалов",
	Factory: New,
}

type updatesCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (u *updatesCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	resp, err := u.f.NewLibrary().GetTvSeriesUpdates(ctx, &emptypb.Empty{}, client.WithRequestTimeout(command.LongRequestTimeout))
	if err != nil {
		u.l.Logf(logger.ErrorLevel, "Get TV-series updates failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	if len(resp.Updates) == 0 {
		return true, command.ReplyText("Новых доступных для загрузки сезонов не найдено")
	}
	var messages []*communication.BotMessage
	for _, r := range resp.Updates {
		messages = append(messages, formatUpdate(r))
	}

	return true, messages
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &updatesCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "updates"}),
	}
}
