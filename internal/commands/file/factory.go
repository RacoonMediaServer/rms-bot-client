package file

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

func newFileHandler(f servicemgr.ServiceFactory, l logger.Logger, mimeType string) command.Command {
	switch mimeType {
	case "application/x-bittorrent":
		return &torrentFileHandler{f: f, l: l}
	default:
		return nil
	}
}
