package checker

import (
	"github.com/wadtech/statusmonitor/service"
	"time"
)

func Listen(work chan *service.Service) {
	for {
		service := <-work
		service.Check()
		work <- service

		time.Sleep(5 * time.Minute)
	}
}
