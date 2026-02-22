package cmdset

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/add"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/archive"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/downloads"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/file"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/library"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/notes"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/remove"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/search"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/snapshot"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/tasks"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/torrents"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/unlink"
)

func New() commands.Factory {
	list := []command.Type{
		archive.Command,
		//cctv.Command,
		downloads.Command,
		file.Command,
		library.Command,
		notes.Command,
		remove.Command,
		search.Command,
		snapshot.Command,
		tasks.AddCommand,
		tasks.ListCommand,
		unlink.Command,
		add.Command,
		torrents.Command,
	}

	cmds := commands.MakeRegisteredCommands(list...)
	return commands.NewDefaultFactory(cmds)
}
