package command

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

// Factory can create Command of specified type. Factory knows all about specific command
type Factory func(interlayer Interlayer, l logger.Logger) Command

type Interlayer struct {
	Services servicemgr.ServiceFactory
}
