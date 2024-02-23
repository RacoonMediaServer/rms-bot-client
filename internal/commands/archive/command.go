package archive

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
	"time"
)

var Command command.Type = command.Type{
	ID:      "archive",
	Title:   "Выгрузка архива",
	Help:    "Выгрузить видео с камеры из архива",
	Factory: New,
}

type archiveCommand struct {
	f       servicemgr.ServiceFactory
	l       logger.Logger
	state   state
	fn      map[state]command.Handler
	cameras map[string]uint32

	camera uint32
	ts     time.Time
	dur    uint64
}

func (c *archiveCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	return c.fn[c.state](ctx)
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	c := &archiveCommand{
		f:       interlayer.Services,
		l:       l.Fields(map[string]interface{}{"command": "archive"}),
		cameras: make(map[string]uint32),
	}

	c.fn = map[state]command.Handler{
		stateInitial:        c.doInitial,
		stateChooseCamera:   c.doChooseCamera,
		stateChooseDay:      c.doChooseDay,
		stateChooseTime:     c.doChooseTime,
		stateChooseDuration: c.doChooseDuration,
	}

	return c
}
