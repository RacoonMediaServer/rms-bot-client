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
	f     servicemgr.ServiceFactory
	l     logger.Logger
	state state
	id    string
}

func (c *tasksListCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch c.state {
	case stateInitial:
		return c.stateInitial(ctx)
	case stateWaitSnoozeDate:
		return c.stateWaitSnoozeDate(ctx)
	default:
		return true, command.ReplyText(command.SomethingWentWrong)
	}
}

func (c *tasksListCommand) stateInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		_, err := c.f.NewNotes().SendTasksNotification(ctx, &rms_notes.SendTasksNotificationRequest{User: ctx.UserID}, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			c.l.Logf(logger.ErrorLevel, "Send tasks notifications failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}
		return true, []*communication.BotMessage{}
	}

	switch ctx.Arguments[0] {
	case "snooze":
		return c.handleSnoozeCommand(ctx)
	case "remove":
		return c.handleRemoveCommand(ctx)
	case "done":
		return c.handleDoneCommand(ctx)
	}

	return true, command.ReplyText(command.ParseArgumentsFailed)
}

func (c *tasksListCommand) handleSnoozeCommand(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	c.id = ctx.Arguments[1]
	c.state = stateWaitSnoozeDate

	return false, []*communication.BotMessage{pickSnoozeDateMessage}
}

func (c *tasksListCommand) handleDoneCommand(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	req := rms_notes.DoneTaskRequest{Id: ctx.Arguments[1], User: ctx.UserID}
	_, err := c.f.NewNotes().DoneTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Done task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача завершена")
}

func (c *tasksListCommand) handleRemoveCommand(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	req := rms_notes.RemoveTaskRequest{Id: ctx.Arguments[1], User: ctx.UserID}
	_, err := c.f.NewNotes().RemoveTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Remove task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача удалена")
}

func (c *tasksListCommand) stateWaitSnoozeDate(ctx command.Context) (bool, []*communication.BotMessage) {
	date, err := parseSnoozeDate(ctx.Arguments.String())
	if err != nil {
		return false, command.ReplyText("Не удалось распарсить дату")
	}
	dateString := date.Format(obsidianDateFormat)

	req := rms_notes.SnoozeTaskRequest{Id: c.id, DueDate: &dateString, User: ctx.UserID}

	_, err = c.f.NewNotes().SnoozeTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Snooze task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача отложена")
}

func NewListCommand(interlayer command.Interlayer, l logger.Logger) command.Command {
	tc := &tasksListCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "tasks"}),
	}

	return middleware.NewNotesAuthCommand(interlayer, l, tc)
}
