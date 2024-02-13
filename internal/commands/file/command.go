package file

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

import (
	"context"
)

var Command command.Type = command.Type{
	ID:       "file",
	Title:    "Загрузить",
	Help:     "Загрузка файлов на сервер",
	Factory:  New,
	Internal: true,
}

type fileCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
	h fileHandler
}

func (c *fileCommand) Do(ctx context.Context, arguments command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	if c.h == nil {
		if attachment == nil {
			return true, command.ReplyText(command.SomethingWentWrong)
		}
		c.h = newFileHandler(c.f, c.l, attachment.MimeType)
		if c.h == nil {
			return true, command.ReplyText("Формат файла не поддерживается")
		}
	}
	return c.h.Do(ctx, arguments, attachment)
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &fileCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "file"}),
	}
}
