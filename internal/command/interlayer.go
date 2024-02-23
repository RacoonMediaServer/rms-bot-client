package command

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/background"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type TaskService interface {
	StartTask(task background.Task)
}

type Interlayer struct {
	Services    servicemgr.ServiceFactory
	TaskService TaskService
}
