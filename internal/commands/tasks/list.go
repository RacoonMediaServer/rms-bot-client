package tasks

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/middleware"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

var ListCommand command.Type = command.Type{
	ID:      "tasks",
	Title:   "Задачи",
	Help:    "Показать текущие",
	Factory: NewListCommand,
}

type tasksListCommand struct {
	f servicemgr.ServiceFactory
	l logger.Logger
}

func (n *tasksListCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	_, err := n.f.NewNotes().SendTasksNotification(ctx, &rms_notes.SendTasksNotificationRequest{User: ctx.UserID}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Send tasks notifications failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, []*communication.BotMessage{}
}

func NewListCommand(interlayer command.Interlayer, l logger.Logger) command.Command {
	tc := &tasksListCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "tasks"}),
	}

	return middleware.NewNotesAuthCommand(interlayer, l, tc)
}
