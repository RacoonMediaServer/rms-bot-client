package bot

import (
	"context"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
)

// Transport is required methods for server session implementation
type Transport interface {
	Receive() <-chan *communication.UserMessage
	Send() chan<- *communication.BotMessage
}

type SpeechRecognizer interface {
	Recognize(ctx context.Context, msg *communication.UserMessage) (string, error)
}

type Settings struct {
	Transport        Transport
	Interlayer       command.Interlayer
	CmdFactory       commands.Factory
	SpeechRecognizer SpeechRecognizer
}
