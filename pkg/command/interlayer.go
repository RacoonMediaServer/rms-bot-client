package command

import (
	"context"
	"reflect"

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
	Extras         map[reflect.Type]any
}

// TODO: very ugly, but very useful
func InterlayerStore(i *Interlayer, value any) {
	if i.Extras == nil {
		i.Extras = map[reflect.Type]any{}
	}
	i.Extras[reflect.TypeOf(value)] = value
}

func InterlayerLoad[T any](i *Interlayer) (T, bool) {
	var defaultVal T
	v, ok := i.Extras[reflect.TypeOf(defaultVal)]
	if !ok {
		return defaultVal, ok
	}

	result, ok := v.(T)
	return result, ok
}
