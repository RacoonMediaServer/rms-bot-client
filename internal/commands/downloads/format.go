package downloads

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"time"
)

func formatTorrent(t *rms_torrent.TorrentInfo) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("<b>%s</b>\n\n", t.Title)
	msg.Text += fmt.Sprintf("<b>Статус</b>: <i>%s</i>", statusToString(t.Status))
	if t.Status == rms_torrent.Status_Downloading {
		msg.Text += fmt.Sprintf("\n<b>Прогресс</b>: %0.2f %%\n", t.Progress)
		msg.Text += fmt.Sprintf("<b>Примерно осталось</b>: %s\n", time.Duration(t.Estimate))
	}
	msg.KeyboardStyle = communication.KeyboardStyle_Message

	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Удалить",
		Command: "/downloads remove " + t.Id,
	})

	if t.Status != rms_torrent.Status_Downloading {
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   "Повысить приоритет",
			Command: "/downloads up " + t.Id,
		})
	}

	return &msg
}

func statusToString(status rms_torrent.Status) string {
	switch status {
	case rms_torrent.Status_Done:
		return "Завершено"
	case rms_torrent.Status_Downloading:
		return "Загружается"
	case rms_torrent.Status_Pending:
		return "В очереди"
	case rms_torrent.Status_Failed:
		return "Ошибка"
	default:
		return ""
	}
}
