package bot

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/logger"
)

func (bot *Bot) process() {
	for {
		select {
		case <-bot.ctx.Done():
			bot.l.Log(logger.DebugLevel, "Shutdown...")
			return

		case msg := <-bot.srv.Receive():
			bot.incomingMessage(msg)
		}
	}
}

func (bot *Bot) incomingMessage(msg *communication.UserMessage) {
	if msg.Type == communication.MessageType_LinkageCode {
		bot.code <- msg.Text
		return
	}
	chat, ok := bot.chats[msg.User]
	if !ok {
		fn := func(m *communication.BotMessage) {
			if m.User == command.BroadcastMessage {
				m.User = 0
			} else {
				m.User = msg.User
			}
			bot.srv.Send() <- m
		}
		chat = newChat(msg.User, bot.interlayer, fn)
		chat.voiceRecognition = bot.voiceRecognition
		bot.chats[msg.User] = chat
	}

	chat.processMessage(msg)
}
