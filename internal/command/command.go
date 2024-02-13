package command

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"strings"
	"time"
)

// LongRequestTimeout is a default timeout for long-possible requests
const LongRequestTimeout = 1 * time.Minute

// Command represents chat bot command
type Command interface {
	// Do executes the command and returns done state and messages to response
	Do(ctx context.Context, arguments Arguments, attachment *communication.Attachment) (done bool, messages []*communication.BotMessage)
}

// IsCommand checks the text can be interpreted as command
func IsCommand(text string) bool {
	if text == "" {
		return false
	}
	return text[0] == '/'
}

// Parse splits text string to command name and arguments
func Parse(text string) (command string, arguments Arguments) {
	list := strings.Split(text, " ")
	if len(list) == 0 {
		return
	}
	command = strings.TrimPrefix(list[0], "/")
	arguments = list[1:]
	return
}
