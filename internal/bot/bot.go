package bot

import (
	"context"
	"sync"

	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server is required methods for server session implementation
type Server interface {
	Receive() <-chan *communication.UserMessage
	Send() chan<- *communication.BotMessage
}

// Bot is chat bot delivery entity
type Bot struct {
	l          logger.Logger
	srv        Server
	interlayer command.Interlayer
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup

	voiceRecognition bool
	chats            map[int32]*chat
	code             chan string
}

// New creates a new chat bot
func New(server Server, interlayer command.Interlayer, voiceRecognitionEnabled bool) *Bot {
	bot := &Bot{
		l:                logger.DefaultLogger.Fields(map[string]interface{}{"from": "bot"}),
		interlayer:       interlayer,
		srv:              server,
		chats:            map[int32]*chat{},
		code:             make(chan string),
		voiceRecognition: voiceRecognitionEnabled,
	}
	bot.interlayer.Messenger = bot
	bot.ctx, bot.cancel = context.WithCancel(context.Background())
	bot.wg.Add(1)
	go func() {
		defer bot.wg.Done()
		bot.process()
	}()
	return bot
}

func (bot *Bot) GetIdentificationCode(ctx context.Context, empty *emptypb.Empty, response *rms_bot_client.GetIdentificationCodeResponse) error {
	bot.l.Log(logger.InfoLevel, "--> GetIdentificationCode()")
	msg := &communication.BotMessage{}
	msg.Type = communication.MessageType_AcquiringCode

	select {
	case bot.srv.Send() <- msg:
	case <-ctx.Done():
		return ctx.Err()
	}

	// TODO: сделать это более надежно
	select {
	case response.Code = <-bot.code:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (bot *Bot) SendMessage(ctx context.Context, req *rms_bot_client.SendMessageRequest, empty *emptypb.Empty) error {
	bot.l.Logf(logger.InfoLevel, "External outgoing message: '%s'", req.Message.Text)
	select {
	case bot.srv.Send() <- req.Message:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (bot *Bot) Shutdown() {
	bot.cancel()
	bot.wg.Wait()
}
