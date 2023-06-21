package tasks

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"time"
)

const obsidianDateFormat = "2006-01-02"

var pickTaskDateMessage *communication.BotMessage = &communication.BotMessage{
	Text:          "Введите дату (yyyy-mm-dd)",
	KeyboardStyle: communication.KeyboardStyle_Chat,
	Buttons: []*communication.Button{
		{
			Command: "Без даты",
		},
		{
			Command: "Завтра",
		},
		{
			Command: "Послезавтра",
		},
	},
}

var pickSnoozeDateMessage *communication.BotMessage = &communication.BotMessage{
	Text:          "На какое число отложить? (yyyy-mm-dd)",
	KeyboardStyle: communication.KeyboardStyle_Chat,
	Buttons: []*communication.Button{
		{
			Command: "Завтра",
		},
		{
			Command: "Послезавтра",
		},
		{
			Command: "На неделю",
		},
		{
			Command: "На месяц",
		},
	},
}

func parseDoneDate(strDate string) (t *time.Time, err error) {
	date := time.Now()
	switch strDate {
	case "Без даты":
		return
	case "Завтра":
		date = date.AddDate(0, 0, 1)
		t = &date
		return
	case "Послезавтра":
		date = date.AddDate(0, 0, 2)
		t = &date
		return
	}
	date, err = time.Parse(obsidianDateFormat, strDate)
	t = &date
	return
}

func parseSnoozeDate(strDate string) (t *time.Time, err error) {
	date := time.Now()
	switch strDate {
	case "Завтра":
		date = date.AddDate(0, 0, 1)
		t = &date
		return
	case "Послезавтра":
		date = date.AddDate(0, 0, 2)
		t = &date
		return
	case "На неделю":
		date = date.AddDate(0, 0, 7)
		t = &date
		return
	case "На месяц":
		date = date.AddDate(0, 1, 0)
		t = &date
		return
	}
	date, err = time.Parse(obsidianDateFormat, strDate)
	t = &date
	return
}
