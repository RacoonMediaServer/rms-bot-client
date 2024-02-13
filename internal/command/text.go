package command

import "github.com/RacoonMediaServer/rms-packages/pkg/communication"

const SomethingWentWrong = "Что-то пошло не так..."
const ParseArgumentsFailed = "Не удалось распознать параметры команды"
const Removed = "Удалено"
const NothingFound = "Ничего не найдено"

func ReplyText(text string) []*communication.BotMessage {
	return []*communication.BotMessage{
		{
			Text: text,
		},
	}
}
