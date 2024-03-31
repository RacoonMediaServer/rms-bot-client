package bot

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-bot-client/internal/commands"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_speech "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-speech"
	"go-micro.dev/v4/logger"
	"time"
)

const voiceRecognitionTimeoutSec = 120

type sendFunc func(msg *communication.BotMessage)

type chat struct {
	l                logger.Logger
	interlayer       command.Interlayer
	send             sendFunc
	voiceRecognition bool

	e *execution
}

func newChat(user int32, interlayer command.Interlayer, send sendFunc) *chat {
	return &chat{
		l:          logger.DefaultLogger.Fields(map[string]interface{}{"chat": user}),
		interlayer: interlayer,
		send:       send,
	}
}

func (c *chat) replyText(text string) {
	c.send(&communication.BotMessage{Text: text})
}

func (c *chat) processMessage(msg *communication.UserMessage) {
	c.l.Logf(logger.InfoLevel, "Got message: %s", msg.Text)
	args := command.Arguments{}

	if c.voiceRecognition && msg.Attachment != nil && msg.Attachment.Type == communication.Attachment_Voice {
		c.recognizeVoice(msg)
		return
	}

	if command.IsCommand(msg.Text) {
		// отменяем предыдущую команду
		if c.e != nil {
			c.e.cancel()
			c.e = nil
		}

		cmdID := ""
		cmdID, args = command.Parse(msg.Text)
		cmd, err := commands.NewCommand(cmdID, c.interlayer, c.l)
		if err != nil {
			c.replyText("Неизвестная команда, всегда можно набрать /help...")
			return
		}
		c.e = newExecution(cmd, c.send, msg.User)

	} else {
		if c.e == nil || c.e.isDone() {
			c.e = nil
			if msg.Attachment != nil {
				c.l.Logf(logger.InfoLevel, "Got file: %s [ %d bytes ]", msg.Attachment.MimeType, len(msg.Attachment.Content))
				cmd, err := commands.NewCommand("file", c.interlayer, c.l)
				if err != nil {
					c.replyText(command.SomethingWentWrong)
					return
				}
				c.e = newExecution(cmd, c.send, msg.User)
			} else {
				c.replyText("Необходимо указать команду. Например: /help")
				return
			}
		}
		args = command.ParseArguments(msg.Text)
	}

	c.e.args <- &execArgs{args: args, attachment: msg.Attachment}
}

func (c *chat) recognizeVoice(msg *communication.UserMessage) {
	speechService := c.interlayer.Services.NewSpeech()
	req := rms_speech.StartRecognitionRequest{
		Data:        msg.Attachment.Content,
		ContentType: msg.Attachment.MimeType,
		TimeoutSec:  voiceRecognitionTimeoutSec,
	}
	resp, err := speechService.StartRecognition(context.TODO(), &req)
	if err != nil {
		c.l.Logf(logger.ErrorLevel, "Start voice recognition failed: %s", err)
		c.replyText("Не удалось распознать голосове сообщение")
		return
	}
	recognized := ""
	for {
		status, err := speechService.GetRecognitionStatus(context.TODO(), &rms_speech.GetRecognitionStatusRequest{JobId: resp.JobId})
		if err != nil {
			c.l.Logf(logger.ErrorLevel, "Get status of voice recognition failed: %s", err)
			c.replyText("Ошибка при попытке распознавания голоса")
			return
		}
		if status.Status == rms_speech.GetRecognitionStatusResponse_Failed {
			c.l.Logf(logger.ErrorLevel, "Voice recognition failed: %s", err)
			c.replyText("Не удалось распознать голосове сообщение")
			return
		}
		if status.Status == rms_speech.GetRecognitionStatusResponse_Done {
			recognized = status.RecognizedText
			break
		}
		<-time.After(1 * time.Second)
	}

	c.replyText(fmt.Sprintf("<b>Распознано</b>: %s", recognized))

	// TODO: execute text command
	if c.e != nil && !c.e.isDone() {
		args := command.ParseArguments(recognized)
		c.e.args <- &execArgs{args: args, attachment: msg.Attachment}
	}
}
