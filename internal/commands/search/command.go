package search

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"time"
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
}

func replyText(text string) []*communication.BotMessage {
	return []*communication.BotMessage{
		{
			Text: text,
		},
	}
}

func (s *searchCommand) Do(ctx context.Context, arguments command.Arguments) (done bool, messages []*communication.BotMessage) {
	if len(arguments) < 1 {
		return false, replyText("Что ищем?")
	}

	resp, err := s.f.NewLibrary().SearchMovie(ctx, &rms_library.SearchMovieRequest{Text: arguments.String(), Limit: searchMoviesLimit}, client.WithRequestTimeout(1*time.Minute))
	if err != nil {
		logger.Errorf("searchMovie failed: %s", err)
		return true, replyText("Что-то пошло не так...")
	}
	logger.Infof("Got %d results", len(resp.Movies))

	var result []*communication.BotMessage
	for _, mov := range resp.Movies {
		result = append(result, formatMovieMessage(mov))
	}

	return true, result
}

func New(f servicemgr.ServiceFactory) command.Command {
	return &searchCommand{f: f}
}
