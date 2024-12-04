package command

import (
	"context"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/background"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskService interface {
	StartTask(task background.Task)
}

type MessageSender interface {
	SendMessage(ctx context.Context, request *rms_bot_client.SendMessageRequest, response *emptypb.Empty) error
}

type Interlayer struct {
	Services       servicemgr.ServiceFactory
	TaskService    TaskService
	Messenger      MessageSender
	ShareDirectory string
}
