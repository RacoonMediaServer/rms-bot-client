package library

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

func formatMovie(mov *rms_library.Movie) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("<b>%s (%d)</b>", mov.Info.Title, mov.Info.Year)

	if mov.TvSeries != nil {
		msg.Text += fmt.Sprintf("\nСкачано сезонов: %d", len(mov.TvSeries.Seasons))
	}
	msg.Text += fmt.Sprintf("\nЗанимаемое место: %.02f Гб", float64(mov.Size)/float64(1024))
	msg.KeyboardStyle = communication.KeyboardStyle_Message
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Удалить",
		Command: "/remove " + mov.Id,
	})
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Докачать",
		Command: "/download watchlist " + mov.Id,
	})

	return &msg
}
