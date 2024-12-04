package background

import (
	"context"
	"go-micro.dev/v4/logger"
	"sync"
)

type Service struct {
	l      logger.Logger
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewService() *Service {
	s := Service{
		l:  logger.DefaultLogger.Fields(map[string]interface{}{"from": "background"}),
		wg: sync.WaitGroup{},
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return &s
}

func (s *Service) StartTask(task Task) {
	s.wg.Add(1)
	go func(task Task) {
		defer s.wg.Done()
		l := s.l.Fields(map[string]interface{}{"task": task.Info()})
		l.Logf(logger.InfoLevel, "Start task")
		if err := task.Run(s.ctx, l); err != nil {
			l.Logf(logger.ErrorLevel, "Task done with error: %s", err)
		} else {
			l.Logf(logger.InfoLevel, "Task done")
		}
	}(task)
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
}
