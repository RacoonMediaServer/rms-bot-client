package commands

import (
	"errors"

	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/archive"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/cctv"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/download"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/downloads"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/file"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/library"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/notes"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/remove"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/search"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/snapshot"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/tasks"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/unlink"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands/updates"
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
	commandMap[tasks.AddCommand.ID] = tasks.AddCommand
	commandMap[tasks.ListCommand.ID] = tasks.ListCommand
	commandMap[unlink.Command.ID] = unlink.Command
	commandMap[snapshot.Command.ID] = snapshot.Command
	commandMap[file.Command.ID] = file.Command
	commandMap[archive.Command.ID] = archive.Command
	commandMap[cctv.Command.ID] = cctv.Command
}

func NewCommand(commandID string, interlayer command.Interlayer, l logger.Logger) (command.Command, error) {
	cmd, ok := commandMap[commandID]
	if !ok {
		return nil, ErrCommandNotFound
	}
	return cmd.Factory(interlayer, l), nil
}
