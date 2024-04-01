package chat

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"time"
)

const commandExecutionTimeout = 10 * time.Minute
const maxMessages = 10

type commandExecution struct {
	chatCtx Context
	cmd     command.Command

	ctx      context.Context
	cancel   context.CancelFunc
	messages chan *communication.UserMessage
}

func newCommandExecutionContext(chatCtx Context, cmd command.Command) *commandExecution {
	e := commandExecution{
		chatCtx:  chatCtx,
		cmd:      cmd,
		messages: make(chan *communication.UserMessage, maxMessages),
	}

	e.ctx, e.cancel = context.WithTimeout(chatCtx.Ctx, commandExecutionTimeout)
	go e.execute()

	return &e
}

func (e *commandExecution) execute() {
	defer e.chatCtx.WaitGroup.Done()
	defer e.cancel()
	for {
		select {
		case <-e.messages:
		case <-e.ctx.Done():
			return
		}
	}
}
