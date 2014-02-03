package service

import (
	"log"
	"net"
)

type Service struct {
	Description string `json:description`
	Host        string `json:host`
	Port        string `json:port`
	Ok          bool
}

func NewService(desc string, host string, port string) *Service {
	return &Service{desc, host, port, false}
}

func (s *Service) Check() {
	con, err := net.Dial("tcp", net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		log.Println(err.Error())
		s.Ok = false
		return
	}

	err = con.Close()
	if err != nil {
		log.Println(err.Error())
		s.Ok = false
		return
	}

	s.Ok = true
}
