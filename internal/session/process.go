package session

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go-micro.dev/v4/logger"
	"net/http"
	"sync"
	"time"
)

const reconnectTimeout = 10 * time.Second

func (s *Session) connect() error {
	h := make(http.Header)
	h.Add("X-Token", s.device)

	var err error
	s.conn, _, err = websocket.DefaultDialer.Dial(s.u.String(), h)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) send(msg *communication.BotMessage) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err = s.conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
		return err
	}

	return nil
}

func (s *Session) process() {
	for {

		if err := s.connect(); err != nil {
			s.l.Logf(logger.ErrorLevel, "Connect to the server failed: %s, try reconnect...", err)
			select {
			case <-time.After(reconnectTimeout):
				continue
			case <-s.ctx.Done():
				return
			}
		}

		s.l.Logf(logger.InfoLevel, "Connected to the server")

		ctx, cancel := context.WithCancel(context.Background())

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.readProcess(cancel)
		}()

	handleMessages:
		for {
			select {
			case msg := <-s.outgoing:
				if err := s.send(msg); err != nil {
					s.l.Logf(logger.ErrorLevel, "Send message failed: %s", err)
					break handleMessages
				}
			case <-s.ctx.Done():
				s.l.Logf(logger.DebugLevel, "Shutdowning...")
				_ = s.conn.Close()
				return
			case <-ctx.Done():
				break handleMessages
			}
		}

		_ = s.conn.Close()
		wg.Wait()
	}
}

func (s *Session) readProcess(cancel context.CancelFunc) {
	defer cancel()
	for {
		_, buf, err := s.conn.ReadMessage()
		if err != nil {
			s.l.Logf(logger.ErrorLevel, "Pick message failed: %s", err)
			return
		}

		var msg communication.UserMessage
		if err = proto.Unmarshal(buf, &msg); err != nil {
			s.l.Logf(logger.ErrorLevel, "Deserialize message failed: %s", err)
		}

		s.incoming <- &msg
	}
}
