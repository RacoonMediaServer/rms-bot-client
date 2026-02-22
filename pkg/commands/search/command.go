package search

import (
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const searchMoviesLimit uint32 = 5

var Command command.Type = command.Type{
	ID:      "search",
	Title:   "Поиск фильмов",
	Help:    "Позволяет искать информацию о фильмах/сериалах и перейти к скачиванию",
	Factory: New,
}

type searchCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (s *searchCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	if len(ctx.Arguments) < 1 {
		return false, command.ReplyText("Что ищем?")
	}

	req := rms_library.MoviesSearchRequest{
		Text:  ctx.Arguments.String(),
		Limit: searchMoviesLimit,
	}

	resp, err := s.f.NewMovies().Search(ctx, &req, client.WithRequestTimeout(1*time.Minute))
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "SearchMovie failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	s.l.Logf(logger.InfoLevel, "Got %d results", len(resp.Movies))

	if len(resp.Movies) == 0 {
		return true, command.ReplyText(command.NothingFound)
	}

	// выводим в обратном порядке,чтобы не мотать ленту в тг
	result := make([]*communication.BotMessage, len(resp.Movies))
	for i, mov := range resp.Movies {
		result[len(result)-i-1] = s.formatMovieMessage(mov)
	}

	return false, result
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &searchCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "search"}),
	}
}
