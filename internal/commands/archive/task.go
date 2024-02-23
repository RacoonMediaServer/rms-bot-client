package archive

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/config"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
	"time"
)

const pollInterval = 1 * time.Second

type task struct {
	cli       rms_transcoder.TranscoderService
	messenger command.MessageSender
	ui        uiVideoMessage
	job       string
	user      int32
}

func (t task) Info() string {
	return fmt.Sprintf("wait_video_transcode_%s_%s", t.ui.Camera, t.job)
}

func (t task) Run(ctx context.Context, l logger.Logger) error {
	defer t.done(ctx, l)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
			done, err := t.trySendVideo(ctx)
			if done {
				return err
			}
		}
	}
}

func (t task) trySendVideo(ctx context.Context) (bool, error) {
	resp, err := t.cli.GetJob(ctx, &rms_transcoder.GetJobRequest{JobId: t.job}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return true, err
	}

	switch resp.Status {
	case rms_transcoder.GetJobResponse_Failed:
		msg := command.ReplyText("Произошла проблема при транскодировании запрошенного видео")
		msg[0].User = t.user
		return true, t.messenger.SendMessage(ctx, &rms_bot_client.SendMessageRequest{Message: msg[0]}, &emptypb.Empty{})
	case rms_transcoder.GetJobResponse_Done:
		content, err := os.ReadFile(filepath.Join(config.Config().ContentDirectory, resp.Destination))
		if err != nil {
			return true, err
		}
		msg := formatVideoMessage(t.ui, content)
		msg.User = t.user
		return true, t.messenger.SendMessage(ctx, &rms_bot_client.SendMessageRequest{Message: msg}, &emptypb.Empty{})
	}

	return false, nil
}

func (t task) done(ctx context.Context, l logger.Logger) {
	req := rms_transcoder.CancelJobRequest{
		JobId:       t.job,
		RemoveFiles: true,
	}
	if _, err := t.cli.CancelJob(ctx, &req, client.WithRequestTimeout(requestTimeout)); err != nil {
		l.Logf(logger.WarnLevel, "Cancel transcoding job failed")
	}
}
