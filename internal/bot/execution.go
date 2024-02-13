package bot

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
)

type execArgs struct {
	args       command.Arguments
	attachment *communication.Attachment
}
type execution struct {
	cmd    command.Command
	ctx    context.Context
	cancel context.CancelFunc
	args   chan *execArgs
	fn     sendFunc
}

const maxArgs = 10

func newExecution(cmd command.Command, fn sendFunc) *execution {
	e := &execution{
		cmd:  cmd,
		args: make(chan *execArgs, maxArgs),
		fn:   fn,
	}
	e.ctx, e.cancel = context.WithCancel(context.TODO())

	go e.execute()
	return e
}

func (e *execution) cancelCommand() {
	e.cancel()
	close(e.args)
}

func (e *execution) isDone() bool {
	select {
	case <-e.ctx.Done():
		return true
	default:
		return false
	}
}

func (e *execution) execute() {
	defer e.cancel()

	for {
		select {
		case args := <-e.args:
			done, replies := e.cmd.Do(e.ctx, args.args, args.attachment)
			for _, m := range replies {
				e.fn(m)
			}
			if done {
				return
			}
		case <-e.ctx.Done():
			return
		}
	}
}
