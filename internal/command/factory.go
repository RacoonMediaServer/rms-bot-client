package command

import (
	"go-micro.dev/v4/logger"
)

// Factory can create Command of specified type. Factory knows all about specific command
type Factory func(interlayer Interlayer, l logger.Logger) Command
