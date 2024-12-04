package file

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
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
	h command.Command
}

func (c *fileCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if c.h == nil {
		if ctx.Attachment == nil {
			return true, command.ReplyText(command.SomethingWentWrong)
		}
		c.h = newFileHandler(c.f, c.l, ctx.Attachment.MimeType)
		if c.h == nil {
			return true, command.ReplyText("Формат файла не поддерживается")
		}
	}
	return c.h.Do(ctx)
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &fileCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "file"}),
	}
}
