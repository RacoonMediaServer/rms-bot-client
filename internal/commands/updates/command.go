package updates

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

import (
	"context"
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

func replyText(text string) []*communication.BotMessage {
	return []*communication.BotMessage{
		{
			Text: text,
		},
	}
}

func (u *updatesCommand) Do(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	resp, err := u.f.NewLibrary().GetTvSeriesUpdates(ctx, &emptypb.Empty{}, client.WithRequestTimeout(command.LongRequestTimeout))
	if err != nil {
		u.l.Logf(logger.ErrorLevel, "Get TV-series updates failed: %s", err)
		return true, replyText(command.SomethingWentWrong)
	}
	if len(resp.Updates) == 0 {
		return true, replyText("Новых доступных для загрузки сезонов не найдено")
	}
	var messages []*communication.BotMessage
	for _, r := range resp.Updates {
		messages = append(messages, formatUpdate(r))
	}

	return true, messages
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &updatesCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "updates"}),
	}
}
