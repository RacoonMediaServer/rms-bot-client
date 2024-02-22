package archive

import (
	"context"
	"fmt"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-packages/pkg/video"
	"github.com/google/uuid"
	"go-micro.dev/v4/client"
)

func (c *archiveCommand) start(ctx context.Context) error {
	// 1. Get archive URL
	replyUri, err := c.getReplyUri(ctx)
	if err != nil {
		return fmt.Errorf("get reply uri failed: %s", err)
	}

	// 2. Create transcoder job
	job, err := c.createJob(ctx, replyUri)
	if err != nil {
		return fmt.Errorf("create transcoding job failed: %s", err)
	}

	// 3. Run monitor
	// TODO

	return nil
}

func (c *archiveCommand) getReplyUri(ctx context.Context) (string, error) {
	ts := uint64(c.ts.UTC().Unix())
	req := rms_cctv.GetReplayUriRequest{
		CameraId:  c.camera,
		Transport: video.Transport_RTSP,
		Timestamp: &ts,
	}
	resp, err := c.f.NewCctv().GetReplayUri(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return "", err
	}
	return resp.Uri, nil
}

func (c *archiveCommand) createJob(ctx context.Context, replyUri string) (string, error) {
	req := rms_transcoder.AddJobRequest{
		Profile:      "telegram",
		Source:       replyUri,
		Destination:  uuid.NewString() + ".mp4",
		AutoComplete: false,
	}
	resp, err := c.f.NewTranscoder().AddJob(ctx, &req, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return "", err
	}
	return resp.TaskId, nil
}
