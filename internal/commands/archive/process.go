package archive

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"time"
)

const shiftReplyDuration = 10 * time.Second

func (c *archiveCommand) start(ctx context.Context, user int32) error {
	c.l.Logf(logger.InfoLevel, "Download camera = %s, time = %s, duration = %d sec", c.ui.Camera, c.ui.Time, c.ui.Duration)
	// 1. Get archive URL
	replyUri, err := c.getReplyUri(ctx)
	if err != nil {
		return fmt.Errorf("get reply uri failed: %s", err)
	}

	c.l.Logf(logger.DebugLevel, "Reply URI = %s", replyUri)

	// 2. Create transcoder job
	job, err := c.createJob(ctx, replyUri)
	if err != nil {
		return fmt.Errorf("create transcoding job failed: %s", err)
	}
	c.l.Logf(logger.DebugLevel, "Job = %s", job)

	// 3. Run monitor
	t := task{
		job:       job,
		cli:       c.interlayer.Services.NewTranscoder(),
		messenger: c.interlayer.Messenger,
		user:      user,
		ui:        c.ui,
	}
	c.interlayer.TaskService.StartTask(&t)

	return nil
}

func (c *archiveCommand) getReplyUri(ctx context.Context) (string, error) {
	ts := uint64(c.ts.Add(-shiftReplyDuration).UTC().Unix())
	req := rms_cctv.GetReplayUriRequest{
		CameraId:  c.camera,
		Transport: media.Transport_RTSP,
		Timestamp: &ts,
	}
	resp, err := c.interlayer.Services.NewCctv().GetReplayUri(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return "", err
	}
	return resp.Uri, nil
}

func (c *archiveCommand) createJob(ctx context.Context, replyUri string) (string, error) {
	dur := uint32(c.dur)
	req := rms_transcoder.AddJobRequest{
		Profile:      "telegram",
		Source:       replyUri,
		Destination:  fmt.Sprintf("telegram/%s_%s_%dsec.mp4", c.ui.Camera, c.ui.Time, c.ui.Duration),
		AutoComplete: false,
		Duration:     &dur,
	}
	resp, err := c.interlayer.Services.NewTranscoder().AddJob(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}
