package library

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

func formatMovie(item *rms_library.ListItem, list rms_library.List) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("<b>%s</b>", item.Title)

	msg.Text += fmt.Sprintf("\nЗанимаемое место: %.02f Гб", float64(item.Size)/float64(1024))
	msg.KeyboardStyle = communication.KeyboardStyle_Message

	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Раздачи",
		Command: fmt.Sprintf("/t %s", item.Id),
	})

	if list != rms_library.List_Favourites {
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "В избранное",
			Command: fmt.Sprintf("/remove %s %d", item.Id, rms_library.List_Favourites),
		})
	}

	if list != rms_library.List_WatchList {
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "К просмотру",
			Command: fmt.Sprintf("/remove %s %d", item.Id, rms_library.List_WatchList),
		})
	}

	if list != rms_library.List_Archive {
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "В архив",
			Command: fmt.Sprintf("/remove %s %d", item.Id, rms_library.List_Archive),
		})
	}

	return &msg
}
