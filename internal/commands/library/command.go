package library

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

const searchMoviesLimit uint32 = 5

var Command command.Type = command.Type{
	ID:      "library",
	Title:   "Библиотека",
	Help:    "Можно посмотреть, что было скачано и добавлено",
	Factory: New,
}

type libraryCommand struct {
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

func (s *libraryCommand) Do(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) == 0 {
		msg := communication.BotMessage{Text: "Что ищем?"}
		msg.KeyboardStyle = communication.KeyboardStyle_Chat
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "Фильмы",
			Command: "Фильмы",
		})
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "Сериалы",
			Command: "Сериалы",
		})
		return false, []*communication.BotMessage{&msg}
	}

	var movieType rms_library.MovieType
	switch arguments[0] {
	case "Фильмы":
		movieType = rms_library.MovieType_Film
	case "Сериалы":
		movieType = rms_library.MovieType_TvSeries
	default:
		return false, replyText("Неверная категория")
	}

	resp, err := s.f.NewLibrary().GetMovies(ctx, &rms_library.GetMoviesRequest{Type: &movieType})
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Get movies failed: %s", err)
		return false, replyText("Что-то пошло не так...")
	}

	if len(resp.Result) == 0 {
		return false, replyText("Ничего не найдено")
	}

	var messages []*communication.BotMessage
	for _, r := range resp.Result {
		messages = append(messages, formatMovie(r))
	}
	return false, messages
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &libraryCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "library"}),
	}
}
