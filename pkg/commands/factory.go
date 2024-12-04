package commands

import (
	"errors"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"go-micro.dev/v4/logger"
)

var ErrCommandNotFound = errors.New("unknown command")

type RegisteredCommands map[string]command.Type

type Factory interface {
	NewCommand(commandID string, interlayer command.Interlayer, l logger.Logger) (command.Command, error)
}

type defaultFactory struct {
	cmds RegisteredCommands
}

// NewCommand implements Factory.
func (d *defaultFactory) NewCommand(commandID string, interlayer command.Interlayer, l logger.Logger) (command.Command, error) {
	cmd, ok := d.cmds[commandID]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return cmd.Factory(interlayer, l), nil
}

func NewDefaultFactory(registeredCommands RegisteredCommands) Factory {
	registeredCommands[helpCommandType.ID] = helpCommandType
	return &defaultFactory{cmds: registeredCommands}
}

func MakeRegisteredCommands(commands ...command.Type) RegisteredCommands {
	result := RegisteredCommands{}
	for _, cmd := range commands {
		result[cmd.ID] = cmd
	}
	return result
}
