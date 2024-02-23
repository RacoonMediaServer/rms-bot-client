package bot

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

type sendFunc func(msg *communication.BotMessage)

type chat struct {
	l    logger.Logger
	f    servicemgr.ServiceFactory
	send sendFunc

	e *execution
}

func newChat(user int32, f servicemgr.ServiceFactory, send sendFunc) *chat {
	return &chat{
		l:    logger.DefaultLogger.Fields(map[string]interface{}{"chat": user}),
		f:    f,
		send: send,
	}
}

func (c *chat) replyText(text string) {
	c.send(&communication.BotMessage{Text: text})
}

func (c *chat) processMessage(msg *communication.UserMessage) {
	c.l.Logf(logger.InfoLevel, "Got message: %s", msg.Text)
	args := command.Arguments{}

	interlayer := command.Interlayer{
		Services: c.f,
	}

	if command.IsCommand(msg.Text) {
		// отменяем предыдущую команду
		if c.e != nil {
			c.e.cancel()
			c.e = nil
		}

		cmdID := ""
		cmdID, args = command.Parse(msg.Text)
		cmd, err := commands.NewCommand(cmdID, interlayer, c.l)
		if err != nil {
			c.replyText("Неизвестная команда, всегда можно набрать /help...")
			return
		}
		c.e = newExecution(cmd, c.send, msg.User)

	} else {
		if c.e == nil || c.e.isDone() {
			c.e = nil
			if msg.Attachment != nil {
				c.l.Logf(logger.InfoLevel, "Got file: %s [ %d bytes ]", msg.Attachment.MimeType, len(msg.Attachment.Content))
				cmd, err := commands.NewCommand("file", interlayer, c.l)
				if err != nil {
					c.replyText(command.SomethingWentWrong)
					return
				}
				c.e = newExecution(cmd, c.send, msg.User)
			} else {
				c.replyText("Необходимо указать команду. Например: /help")
				return
			}
		}
		args = command.ParseArguments(msg.Text)
	}

	c.e.args <- &execArgs{args: args, attachment: msg.Attachment}
}
