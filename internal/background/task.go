package background

import (
	"context"
	"go-micro.dev/v4/logger"
)

type Task interface {
	Info() string
	Run(ctx context.Context, l logger.Logger) error
}
