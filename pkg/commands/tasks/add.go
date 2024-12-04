package tasks

import (
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/internal/middleware"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const requestTimeout = 20 * time.Second

var AddCommand command.Type = command.Type{
	ID:      "tasks_add",
	Title:   "Задачи",
	Help:    "Добавить задачу",
	Factory: NewAddCommand,
}

type state int

const (
	stateInitial state = iota
	stateWaitTaskText
	stateWaitTaskDate
	stateWaitSnoozeDate
)

type tasksAddCommand struct {
	f     servicemgr.ServiceFactory
	l     logger.Logger
	title string
	state state
	date  time.Time
}

func (n *tasksAddCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch n.state {
	case stateInitial:
		return n.stateInitial(ctx)
	case stateWaitTaskText:
		return n.stateWaitTaskText(ctx)
	case stateWaitTaskDate:
		return n.stateWaitTaskDate(ctx)
	default:
		return true, command.ReplyText(command.SomethingWentWrong)
	}
}

func (n *tasksAddCommand) stateInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		n.state = stateWaitTaskText
		return false, command.ReplyText("Введите описание задачи")
	}

	n.title = ctx.Arguments.String()
	n.state = stateWaitTaskDate

	return false, []*communication.BotMessage{pickTaskDateMessage}
}

func (n *tasksAddCommand) stateWaitTaskText(ctx command.Context) (bool, []*communication.BotMessage) {
	n.title = ctx.Arguments.String()
	n.state = stateWaitTaskDate

	return false, []*communication.BotMessage{pickTaskDateMessage}
}

func (n *tasksAddCommand) stateWaitTaskDate(ctx command.Context) (bool, []*communication.BotMessage) {
	date, err := parseDoneDate(ctx.Arguments.String())
	if err != nil {
		return false, command.ReplyText("Не удалось распарсить дату")
	}

	req := rms_notes.AddTaskRequest{Text: n.title, User: ctx.UserID}
	if date != nil {
		dateString := date.Format(obsidianDateFormat)
		req.DueDate = &dateString
	}

	_, err = n.f.NewNotes().AddTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Add task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача добавлена")
}

func NewAddCommand(interlayer command.Interlayer, l logger.Logger) command.Command {
	tc := &tasksAddCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "tasks-add"}),
	}

	return middleware.NewNotesAuthCommand(interlayer, l, tc)
}
