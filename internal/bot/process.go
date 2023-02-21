package bot

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			m.Type = communication.MessageType_Interaction
			m.User = msg.User
			m.Timestamp = timestamppb.Now()
			bot.srv.Send() <- m
		}
		chat = newChat(msg.User, bot.f, fn)
		bot.chats[msg.User] = chat
	}

	chat.processMessage(msg)
}
