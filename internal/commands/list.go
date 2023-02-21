package commands

import (
	"errors"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/search"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

var commandMap map[string]command.Type

var ErrCommandNotFound = errors.New("unknown command")

func init() {
	commandMap = map[string]command.Type{}
	commandMap[helpCommandType.ID] = helpCommandType
	commandMap[search.Command.ID] = search.Command
}

func NewCommand(commandID string, f servicemgr.ServiceFactory) (command.Command, error) {
	cmd, ok := commandMap[commandID]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return cmd.Factory(f), nil
}
