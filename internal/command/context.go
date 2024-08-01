package command

import (
	"context"
	"time"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
)

const BroadcastMessage int32 = 1

type Context struct {
	Ctx        context.Context
	Arguments  Arguments
	Attachment *communication.Attachment
	UserID     int32
}

func (ctx Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.Ctx.Deadline()
}

func (ctx Context) Done() <-chan struct{} {
	return ctx.Ctx.Done()
}

func (ctx Context) Err() error {
	return ctx.Ctx.Err()
}

func (ctx Context) Value(key any) any {
	return ctx.Ctx.Value(key)
}
