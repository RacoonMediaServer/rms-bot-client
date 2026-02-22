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
		msg := communication.BotMessage{Text: "Какой список отобразить?"}
		msg.KeyboardStyle = communication.KeyboardStyle_Chat
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   string(listNameFavourites),
			Command: string(listNameFavourites),
		})
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   string(listNameWatchList),
			Command: string(listNameWatchList),
		})
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   string(listNameArchive),
			Command: string(listNameArchive),
		})
		return false, []*communication.BotMessage{&msg}
	}

	listType, ok := getListType(listName(ctx.Arguments[0]))
	if !ok {
		return false, command.ReplyText("Неверная категория")
	}

	movieType := rms_library.ContentType_TypeMovies
	sort := rms_library.Sort{Order: rms_library.Sort_Desc}
	if listType != rms_library.List_Favourites {
		sort.By = rms_library.Sort_CreatedAt
		sort.Order = rms_library.Sort_Asc
	}
	req := rms_library.ListsListRequest{
		List:        listType,
		ContentType: &movieType,
		Sort:        &sort,
	}
	resp, err := s.f.NewLists().List(ctx, &req)
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Get movies failed: %s", err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	if len(resp.Items) == 0 {
		return false, command.ReplyText(command.NothingFound)
	}

	messages := make([]*communication.BotMessage, len(resp.Items))
	for i, r := range resp.Items {
		messages[len(messages)-i-1] = formatMovie(r, listType)
	}
	return false, messages
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &libraryCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "library"}),
	}
}
