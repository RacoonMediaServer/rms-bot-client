package snapshot

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"time"
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

func formatSnapshot(cameraName string, snapshot []byte) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("Изображение с камеры '%s' (%s)", cameraName, time.Now().Format(time.RFC1123))
	msg.Attachment = &communication.Attachment{
		Type:     communication.Attachment_Photo,
		MimeType: "image/jpeg",
		Content:  snapshot,
	}
	return &msg
}
