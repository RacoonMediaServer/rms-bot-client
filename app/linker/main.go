package main

import (
	"context"
	"fmt"
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
}
