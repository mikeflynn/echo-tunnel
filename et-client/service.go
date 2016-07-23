package main

import (
	"log"

	"github.com/kardianos/service"
)

var Logger service.Logger

type Program struct{}

func (p *Program) Start(s service.Service) error {
	go run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	StopFlag = true
	return nil
}
