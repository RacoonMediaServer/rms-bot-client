package main

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	var device string
	service := micro.NewService(
		micro.Name("rms-bot-client.linker"),
		micro.Flags(
			&cli.StringFlag{
				Name:        "device",
				Usage:       "Device ID",
				Required:    true,
				Destination: &device,
			},
		),
	)
	service.Init()

	f := servicemgr.NewServiceFactory(service)
	resp, err := f.NewBotClient().GetIdentificationCode(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Code)

	_, err = f.NewBotClient().SendMessage(context.Background(), &rms_bot_client.SendMessageRequest{Message: &communication.BotMessage{Text: "Identification code requested"}})
	if err != nil {
		panic(err)
	}
}
