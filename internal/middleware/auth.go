package middleware

import (
	"github.com/RacoonMediaServer/rms-bot-client/internal/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"time"
)

type authMiddleware struct {
	f     servicemgr.ServiceFactory
	l     logger.Logger
	cmd   command.Command
	state authState

	savedCtx command.Context
	req      *rms_notes.UserLoginRequest
}

type authState int

const (
	authStateInitial authState = iota
	authStateWaitEndpoint
	authStateWaitLogin
	authStateWaitPassword
	authStatePassThrough
)

const notesReqTimeout = 10 * time.Second

func NewNotesAuthCommand(interlayer command.Interlayer, l logger.Logger, cmd command.Command) command.Command {
	return &authMiddleware{
		f:   interlayer.Services,
		l:   l.Fields(map[string]interface{}{"middleware": "notes-auth"}),
		cmd: cmd,
	}
}

func (a *authMiddleware) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	switch a.state {
	case authStateInitial:
		return a.stateInitial(ctx)
	case authStateWaitEndpoint:
		return a.stateWaitEndpoint(ctx)
	case authStateWaitLogin:
		return a.stateWaitLogin(ctx)
	case authStateWaitPassword:
		return a.stateWaitPassword(ctx)
	case authStatePassThrough:
		return a.cmd.Do(ctx)
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}

func (a *authMiddleware) stateInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	cli := a.f.NewNotes()
	resp, err := cli.IsUserLogged(ctx, &rms_notes.IsUserLoggedRequest{User: ctx.UserID}, client.WithRequestTimeout(notesReqTimeout))
	if err != nil {
		a.l.Logf(logger.ErrorLevel, "Cannot access to notes: %s", err)
		return true, command.ReplyText("Сервис заметок недоступен")
	}
	if resp.Result {
		a.state = authStatePassThrough
		return a.cmd.Do(ctx)
	}

	a.l.Logf(logger.WarnLevel, "User is not logged, login...")

	a.req = &rms_notes.UserLoginRequest{User: ctx.UserID}
	a.savedCtx = ctx

	msg := communication.BotMessage{}
	msg.Text = "Для доступа к заметкам необходимо указать данные учетной записи Nextcloud. Введите путь до сервера, например http://nc.rms.local/remote.php/dav"
	msg.KeyboardStyle = communication.KeyboardStyle_Message
	msg.Buttons = []*communication.Button{
		{Title: "Определить автоматически", Command: "auto"},
	}
	a.state = authStateWaitEndpoint
	return false, []*communication.BotMessage{&msg}
}

func (a *authMiddleware) stateWaitEndpoint(ctx command.Context) (bool, []*communication.BotMessage) {
	endpoint := ctx.Arguments.String()
	if endpoint == "auto" {
		endpoint = "http://nc.rms.local/remote.php/dav"
	}

	a.req.Endpoint = endpoint
	a.state = authStateWaitLogin

	return false, command.ReplyText("Введите имя пользователя")
}

func (a *authMiddleware) stateWaitLogin(ctx command.Context) (bool, []*communication.BotMessage) {
	a.req.Login = ctx.Arguments.String()
	a.state = authStateWaitPassword

	return false, command.ReplyText("Введите пароль")
}

func (a *authMiddleware) stateWaitPassword(ctx command.Context) (bool, []*communication.BotMessage) {
	a.req.Password = ctx.Arguments.String()

	cli := a.f.NewNotes()
	resp, err := cli.UserLogin(ctx, a.req, client.WithRequestTimeout(notesReqTimeout))
	if err != nil {
		a.l.Logf(logger.ErrorLevel, "Login to notes failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	switch resp.Code {
	case rms_notes.UserLoginResponse_ConnectionError:
		return true, command.ReplyText("Не удалось подключиться к Nextcloud")
	case rms_notes.UserLoginResponse_InvalidCredentials:
		return true, command.ReplyText("Неверно указаны данные учетной записи")
	case rms_notes.UserLoginResponse_OK:
		a.state = authStatePassThrough
		reply := command.ReplyText("Сервис заметок подключен")
		result, items := a.cmd.Do(a.savedCtx)
		reply = append(reply, items...)
		return result, reply
	}

	return true, command.ReplyText(command.SomethingWentWrong)
}
