package archive

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

func formatCameraList(cameras []*rms_cctv.Camera) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Text = "Выберите камеру"

	for _, cam := range cameras {
		msg.Buttons = append(msg.Buttons, &communication.Button{
			Title:   cam.Name,
			Command: cam.Name,
		})
	}

	return &msg
}

func formatDayRequest() *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Text = "Введите дату в формате YYYY-MM-DD"

	msg.Buttons = []*communication.Button{
		{
			Title:   "Сегодня",
			Command: "Сегодня",
		},
		{
			Title:   "Вчера",
			Command: "Вчера",
		},
		{
			Title:   "Позавчера",
			Command: "Позавчера",
		},
	}

	return &msg
}

func formatTimeRequest() *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Text = "Введите время в формате hh:mm:ss"

	return &msg
}

func formatDurationRequest() *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Text = "Выберите длительность ролика (секунды)"

	msg.Buttons = []*communication.Button{
		{
			Title:   "30",
			Command: "30",
		},
		{
			Title:   "60",
			Command: "60",
		},
		{
			Title:   "120",
			Command: "120",
		},
		{
			Title:   "240",
			Command: "240",
		},
	}

	return &msg
}
