package chat

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"sync"
)

type Context struct {
	Ctx        context.Context
	Interlayer command.Interlayer
	Send       chan<- *communication.BotMessage
	User       int32
	WaitGroup  *sync.WaitGroup
}

type Chat struct {
	chatCtx Context
}

func New(ctx Context) *Chat {
	return &Chat{
		chatCtx: ctx,
	}
}

func (c *Chat) PushMessage(msg *communication.UserMessage) {

}
