package commands

import (
	"fmt"
	"sort"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/logger"
)

var helpCommandType = command.Type{
	ID:      "help",
	Title:   "Справка",
	Help:    "Пояснить за функции бота",
	Factory: newHelpCommand,
}

type helpCommand struct {
	cmds RegisteredCommands
}

func (h helpCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	titles := make([]string, 0, len(h.cmds))
	for k, _ := range h.cmds {
		titles = append(titles, k)
	}
	sort.Slice(titles, func(i, j int) bool {
		return titles[i] < titles[j]
	})
	result := ""
	for _, t := range titles {
		cmd := h.cmds[t]
		if !cmd.Internal {
			result += fmt.Sprintf("/%s %s - %s\n", cmd.ID, cmd.Title, cmd.Help)
		}
	}
	return true, []*communication.BotMessage{
		{
			Text: result,
		},
	}
}

func newHelpCommand(interlayer command.Interlayer, l logger.Logger) command.Command {
	return &helpCommand{}
}
