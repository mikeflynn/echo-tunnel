package main

import (
	"os"

	"github.com/kardianos/service"
)

type Program struct{}

func (p *Program) Start(s service.Service) error {
	Debug("STARTING the Echo Tunnel client service.")
	go run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	Debug("STOPPING the Echo Tunnel client service.")
	StopFlag = true
	return nil
}

func startService() {
	svcConfig := &service.Config{
		Name:        "EchoTunnelClient",
		DisplayName: "Echo Tunnel Client",
		Description: "Client for the Echo Tunnel service.",
	}

	prgm := &Program{}
	svc, err := service.New(prgm, svcConfig)
	if err != nil {
		Debug(err.Error())
		os.Exit(1)
	}

	err = svc.Run()
	if err != nil {
		Debug(err.Error())
		os.Exit(1)
	}
}
