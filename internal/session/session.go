package session

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-bot-client/internal/config"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/gorilla/websocket"
	"go-micro.dev/v4/logger"
	"net/url"
	"sync"
)

const maxMessagesInQueue = 50

// Session is a connection to the server
type Session struct {
	l      logger.Logger
	u      url.URL
	device string
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	conn *websocket.Conn

	incoming chan *communication.UserMessage
	outgoing chan *communication.BotMessage
}

func New(endpoint config.Remote, deviceID string) *Session {
	s := &Session{
		l: logger.DefaultLogger.Fields(map[string]interface{}{"from": "session"}),
		u: url.URL{
			Scheme: endpoint.Scheme,
			Host:   fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port),
			Path:   endpoint.Path,
		},
		device:   deviceID,
		incoming: make(chan *communication.UserMessage, maxMessagesInQueue),
		outgoing: make(chan *communication.BotMessage, maxMessagesInQueue),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.process()
	}()

	return s
}

func (s *Session) Receive() <-chan *communication.UserMessage {
	return s.incoming
}

func (s *Session) Send() chan<- *communication.BotMessage {
	return s.outgoing
}

func (s *Session) Shutdown() {
	s.cancel()
	s.wg.Wait()
}
