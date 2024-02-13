package file

import (
	"context"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

type fileHandler interface {
	Do(ctx context.Context, args command.Arguments, attachment *communication.Attachment) (bool, []*communication.BotMessage)
}

func newFileHandler(f servicemgr.ServiceFactory, l logger.Logger, mimeType string) fileHandler {
	switch mimeType {
	case "application/x-bittorrent":
		return &torrentFileHandler{f: f, l: l}
	default:
		return nil
	}
}
