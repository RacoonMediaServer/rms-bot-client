package tasks

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"time"
)

const obsidianDateFormat = "2006-01-02"

func composePickTaskDateMessage() *communication.BotMessage {
	now := time.Now()
	return &communication.BotMessage{
		Text:          "Введите дату (yyyy-mm-dd)",
		KeyboardStyle: communication.KeyboardStyle_Chat,
		Buttons: []*communication.Button{
			{
				Title:   "Без даты",
				Command: "0",
			},
			{
				Title:   "Завтра",
				Command: now.AddDate(0, 0, 1).Format(obsidianDateFormat),
			},
			{
				Title:   "Послезавтра",
				Command: now.AddDate(0, 0, 2).Format(obsidianDateFormat),
			},
		},
	}
}

func composePickSnoozeDateMessage() *communication.BotMessage {
	now := time.Now()
	return &communication.BotMessage{
		Text:          "На какое число отложить? (yyyy-mm-dd)",
		KeyboardStyle: communication.KeyboardStyle_Chat,
		Buttons: []*communication.Button{
			{
				Title:   "Завтра",
				Command: now.AddDate(0, 0, 1).Format(obsidianDateFormat),
			},
			{
				Title:   "Послезавтра",
				Command: now.AddDate(0, 0, 2).Format(obsidianDateFormat),
			},
			{
				Title:   "На неделю",
				Command: now.AddDate(0, 0, 7).Format(obsidianDateFormat),
			},
			{
				Title:   "На месяц",
				Command: now.AddDate(0, 1, 0).Format(obsidianDateFormat),
			},
		},
	}
}

func parseDate(strDate string) (*time.Time, error) {
	if strDate == "0" {
		return nil, nil
	}
	date, err := time.Parse(obsidianDateFormat, strDate)
	if err != nil {
		return nil, err
	}

	return &date, nil
}
