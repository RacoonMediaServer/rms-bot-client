package commands

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/logger"
	"sort"
)

var helpCommandType = command.Type{
	ID:      "help",
	Title:   "Справка",
	Help:    "Пояснить за функции бота",
	Factory: newHelpCommand,
}

type helpCommand struct {
}

func (h helpCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	titles := make([]string, 0, len(commandMap))
	for k, _ := range commandMap {
		titles = append(titles, k)
	}
	sort.Slice(titles, func(i, j int) bool {
		return titles[i] < titles[j]
	})
	result := ""
	for _, t := range titles {
		cmd := commandMap[t]
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
