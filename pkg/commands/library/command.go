package library

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
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

func (s *libraryCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
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
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "Ролики",
			Command: "Ролики",
		})
		return false, []*communication.BotMessage{&msg}
	}

	var movieType rms_library.MovieType
	switch ctx.Arguments[0] {
	case "Фильмы":
		movieType = rms_library.MovieType_Film
	case "Сериалы":
		movieType = rms_library.MovieType_TvSeries
	case "Ролики":
		movieType = rms_library.MovieType_Clip
	default:
		return false, command.ReplyText("Неверная категория")
	}

	resp, err := s.f.NewMovies().List(ctx, &rms_library.GetMoviesRequest{Type: &movieType})
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Get movies failed: %s", err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	if len(resp.Result) == 0 {
		return false, command.ReplyText(command.NothingFound)
	}

	messages := make([]*communication.BotMessage, len(resp.Result))
	for i, r := range resp.Result {
		messages[len(messages)-i-1] = formatMovie(r)
	}
	return false, messages
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &libraryCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "library"}),
	}
}
