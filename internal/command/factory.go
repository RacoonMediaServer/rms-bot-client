package command

import "github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"

// Factory can create Command of specified type. Factory knows all about specific command
type Factory func(f servicemgr.ServiceFactory) Command
