package notes

import (
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/middleware"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
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

func (n *notesCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch n.state {
	case stateInitial:
		return n.stateInitial(ctx)
	case stateWaitTitle:
		return n.stateWaitTitle(ctx)
	case stateWaitText:
		return n.stateWaitText(ctx)
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}

func (n *notesCommand) stateInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		n.state = stateWaitTitle
		return false, command.ReplyText("Введите заголовок заметки")
	}

	n.title = ctx.Arguments.String()
	n.state = stateWaitText
	return false, command.ReplyText("Введите текст заметки")
}

func (n *notesCommand) stateWaitTitle(ctx command.Context) (bool, []*communication.BotMessage) {
	n.title = ctx.Arguments.String()
	n.state = stateWaitText
	return false, command.ReplyText("Введите текст заметки")
}

func (n *notesCommand) stateWaitText(ctx command.Context) (bool, []*communication.BotMessage) {
	req := rms_notes.AddNoteRequest{
		Title: n.title,
		Text:  ctx.Arguments.String(),
		User:  ctx.UserID,
	}
	_, err := n.f.NewNotes().AddNote(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		n.l.Logf(logger.ErrorLevel, "Add note failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}
	return true, command.ReplyText("Заметка добавлена")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	nc := &notesCommand{
		f: interlayer.Services,
		l: l.Fields(map[string]interface{}{"command": "notes"}),
	}

	return middleware.NewNotesAuthCommand(interlayer, l, nc)
}
