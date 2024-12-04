package speech

import (
	"context"
	"time"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_speech "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-speech"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

const voiceRecognitionTimeoutSec = 120
const statusRequestInterval = 1 * time.Second

type Recognizer struct {
	ServiceFactory servicemgr.ServiceFactory
}

func (r Recognizer) Recognize(ctx context.Context, msg *communication.UserMessage) (string, error) {
	speechService := r.ServiceFactory.NewSpeech()
	req := rms_speech.StartRecognitionRequest{
		Data:        msg.Attachment.Content,
		ContentType: msg.Attachment.MimeType,
		TimeoutSec:  voiceRecognitionTimeoutSec,
	}

	resp, err := speechService.StartRecognition(ctx, &req)
	if err != nil {
		return "", err
	}
	recognized := ""
	for {
		status, err := speechService.GetRecognitionStatus(ctx, &rms_speech.GetRecognitionStatusRequest{JobId: resp.JobId})
		if err != nil {
			return "", err
		}
		if status.Status == rms_speech.GetRecognitionStatusResponse_Failed {
			return "", err
		}
		if status.Status == rms_speech.GetRecognitionStatusResponse_Done {
			recognized = status.RecognizedText
			break
		}
		select {
		case <-time.After(statusRequestInterval):
		case <-ctx.Done():
			_, _ = speechService.StopRecognition(context.Background(), &rms_speech.StopRecognitionRequest{JobId: resp.JobId})
			return "", ctx.Err()
		}
	}

	return recognized, nil
}
