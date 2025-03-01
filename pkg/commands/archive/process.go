package archive

import (
	"context"
	"fmt"
	"time"

	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const shiftTime = 5 * time.Second

func (c *archiveCommand) start(ctx context.Context, user int32) error {
	c.l.Logf(logger.InfoLevel, "Download camera = %s, time = %s, duration = %d sec", c.ui.Camera, c.ui.Time, c.ui.Duration)
	// 1. Get archive URL
	replyUri, offset, err := c.getReplyUri(ctx)
	if err != nil {
		return fmt.Errorf("get reply uri failed: %s", err)
	}
	c.l.Logf(logger.DebugLevel, "Reply URI = %s, offset = %d", replyUri, offset)

	// 2. Create transcoder job
	job, err := c.createJob(ctx, replyUri, offset)
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
		dir:       c.interlayer.ShareDirectory,
	}
	c.interlayer.TaskService.StartTask(&t)

	return nil
}

func (c *archiveCommand) getReplyUri(ctx context.Context) (uri string, offset uint32, err error) {
	ts := uint64(c.ts.Add(-shiftTime).UTC().Unix())
	req := rms_cctv.GetReplayUriRequest{
		CameraId:  c.camera,
		Transport: media.Transport_RTSP,
		Timestamp: &ts,
	}
	var resp *rms_cctv.GetReplayUriResponse
	resp, err = c.interlayer.Services.NewCctvCameras().GetReplayUri(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return
	}
	uri = resp.Uri
	offset = resp.OffsetSec
	return
}

func (c *archiveCommand) createJob(ctx context.Context, replyUri string, offset uint32) (string, error) {
	dur := uint32(c.dur)
	req := rms_transcoder.AddJobRequest{
		Profile:      "telegram",
		Source:       replyUri,
		Destination:  fmt.Sprintf("telegram/%s_%s_%dsec.mp4", c.ui.Camera, c.ui.Time, c.ui.Duration),
		AutoComplete: false,
		Duration:     &dur,
		Offset:       &offset,
	}
	resp, err := c.interlayer.Services.NewTranscoder().AddJob(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}
