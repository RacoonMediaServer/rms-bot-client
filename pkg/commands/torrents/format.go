package torrents

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

func newAddMessage(id string) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = "Список раздач"
	msg.KeyboardStyle = communication.KeyboardStyle_Message

	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Добавить",
		Command: fmt.Sprintf("/torrents %s add", id),
	})
	return &msg
}

func newTorrentMessage(id string, torrent *rms_library.Torrent) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("%s [ %.2f Gb, %d seeds]", torrent.Title, float32(torrent.Size)/1024., torrent.Seeders)
	msg.KeyboardStyle = communication.KeyboardStyle_Message

	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Удалить",
		Command: fmt.Sprintf("/torrents %s remove %s", id, torrent.Id),
	})
	return &msg
}

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

func formatTorrents(torrents []*rms_library.Torrent) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	for i, t := range torrents {
		msg.Text += fmt.Sprintf("%d. %s [ %.2f Gb, %d seeds]\n", i+1, t.Title, float32(t.Size)/1024., t.Seeders)
		no := fmt.Sprintf("%d", i+1)
		msg.Buttons = append(msg.Buttons, &communication.Button{Title: no, Command: no})
	}
	return &msg
}
