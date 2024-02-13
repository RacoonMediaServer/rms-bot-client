package tasks

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"time"
)

const requestTimeout = 20 * time.Second

var Command command.Type = command.Type{
	ID:      "tasks",
	Title:   "Задачи",
	Help:    "Добавить задачу",
	Factory: New,
}

type state int

const (
	stateInitial state = iota
	stateWaitTaskText
	stateWaitTaskDate
	stateWaitSnoozeDate
)

type tasksCommand struct {
	f     servicemgr.ServiceFactory
	l     logger.Logger
	title string
	id    string
	state state
	date  time.Time
}

func (n *tasksCommand) Do(ctx context.Context, arguments command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	switch n.state {
	case stateInitial:
		return n.stateInitial(ctx, arguments)
	case stateWaitTaskText:
		return n.stateWaitTaskText(ctx, arguments)
	case stateWaitTaskDate:
		return n.stateWaitTaskDate(ctx, arguments)
	case stateWaitSnoozeDate:
		return n.stateWaitSnoozeDate(ctx, arguments)

	}

	return true, command.ReplyText(command.SomethingWentWrong)
}

func (n *tasksCommand) stateInitial(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) == 0 {
		n.state = stateWaitTaskText
		return false, command.ReplyText("Введите описание задачи")
	}

	switch arguments[0] {
	case "snooze":
		return n.handleSnoozeCommand(ctx, arguments)
	case "done":
		return n.handleDoneCommand(ctx, arguments)
	}

	return true, command.ReplyText(command.ParseArgumentsFailed)
}

func (n *tasksCommand) handleSnoozeCommand(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	n.id = arguments[1]
	n.state = stateWaitSnoozeDate

	return false, []*communication.BotMessage{pickSnoozeDateMessage}
}

func (n *tasksCommand) handleDoneCommand(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	req := rms_notes.DoneTaskRequest{Id: arguments[1]}
	_, err := n.f.NewNotes().DoneTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Done task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача завершена")
}

func (n *tasksCommand) stateWaitTaskText(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	n.title = arguments.String()
	n.state = stateWaitTaskDate

	return false, []*communication.BotMessage{pickTaskDateMessage}
}

func (n *tasksCommand) stateWaitTaskDate(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	date, err := parseDoneDate(arguments.String())
	if err != nil {
		return false, command.ReplyText("Не удалось распарсить дату")
	}

	req := rms_notes.AddTaskRequest{Text: n.title}
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

func (n *tasksCommand) stateWaitSnoozeDate(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	date, err := parseSnoozeDate(arguments.String())
	if err != nil {
		return false, command.ReplyText("Не удалось распарсить дату")
	}
	dateString := date.Format(obsidianDateFormat)

	req := rms_notes.SnoozeTaskRequest{Id: n.id, DueDate: &dateString}

	_, err = n.f.NewNotes().SnoozeTask(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Snooze task failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Задача отложена")
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &tasksCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "tasks"}),
	}
}
