package updates

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

func formatSeasons(seasons []uint32) string {
	result := ""
	for _, s := range seasons {
		result += fmt.Sprintf("%d, ", s)
	}
	if len(result) > 2 {
		result = result[0 : len(result)-2]
	}
	return result
}

func formatUpdate(u *rms_library.TvSeriesUpdate) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Message

	msg.Text = fmt.Sprintf("<b>%s (%d)</b>\n", u.Info.Title, u.Info.Year)
	msg.Text += fmt.Sprintf("Доступные сезоны: %s", formatSeasons(u.SeasonsAvailable))

	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Скачать",
		Command: "/download auto " + u.Id,
	})
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Выбрать раздачу",
		Command: "/download select" + u.Id,
	})

	return &msg
}
