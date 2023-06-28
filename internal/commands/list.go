package commands

import (
	"errors"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/download"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/downloads"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/library"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/notes"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/remove"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/search"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/tasks"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/unlink"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/updates"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

var commandMap map[string]command.Type

var ErrCommandNotFound = errors.New("unknown command")

func init() {
	commandMap = map[string]command.Type{}
	commandMap[helpCommandType.ID] = helpCommandType
	commandMap[search.Command.ID] = search.Command
	commandMap[download.Command.ID] = download.Command
	commandMap[library.Command.ID] = library.Command
	commandMap[remove.Command.ID] = remove.Command
	commandMap[downloads.Command.ID] = downloads.Command
	commandMap[updates.Command.ID] = updates.Command
	commandMap[notes.Command.ID] = notes.Command
	commandMap[tasks.Command.ID] = tasks.Command
	commandMap[unlink.Command.ID] = unlink.Command
}

func NewCommand(commandID string, f servicemgr.ServiceFactory, l logger.Logger) (command.Command, error) {
	cmd, ok := commandMap[commandID]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return cmd.Factory(f, l), nil
}
