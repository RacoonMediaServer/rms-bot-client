package notes

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
	ID:      "notes",
	Title:   "Заметки",
	Help:    "Добавить заметку",
	Factory: New,
}

type state int

const (
	stateInitial state = iota
	stateWaitTitle
	stateWaitText
)

type notesCommand struct {
	f     servicemgr.ServiceFactory
	l     logger.Logger
	title string
	state state
}

func (n *notesCommand) Do(ctx context.Context, arguments command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage) {
	switch n.state {
	case stateInitial:
		return n.stateInitial(ctx, arguments)
	case stateWaitTitle:
		return n.stateWaitTitle(ctx, arguments)
	case stateWaitText:
		return n.stateWaitText(ctx, arguments)
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}

func (n *notesCommand) stateInitial(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	if len(arguments) == 0 {
		n.state = stateWaitTitle
		return false, command.ReplyText("Введите заголовок заметки")
	}

	n.title = arguments.String()
	n.state = stateWaitText
	return false, command.ReplyText("Введите текст заметки")
}

func (n *notesCommand) stateWaitTitle(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	n.title = arguments.String()
	n.state = stateWaitText
	return false, command.ReplyText("Введите текст заметки")
}

func (n *notesCommand) stateWaitText(ctx context.Context, arguments command.Arguments) (bool, []*communication.BotMessage) {
	req := rms_notes.AddNoteRequest{
		Title: n.title,
		Text:  arguments.String(),
	}
	_, err := n.f.NewNotes().AddNote(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Add note failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText("Заметка добавлена")
}

func New(f servicemgr.ServiceFactory, l logger.Logger) command.Command {
	return &notesCommand{
		f: f,
		l: l.Fields(map[string]interface{}{"command": "notes"}),
	}
}
