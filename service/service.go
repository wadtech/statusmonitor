package service

import (
	"log"
	"net"
)

type Service struct {
	Description string
	host        string
	port        string
	Ok          bool
}

func NewService(desc string, host string, port string) *Service {
	return &Service{desc, host, port, false}
}

func (s *Service) Check() {
	con, err := net.Dial("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		log.Println(s.Description, "could not connect to", s.host, s.port, "with message", err)
		s.Ok = false
		return
	}

	err = con.Close()
	if err != nil {
		log.Println(s.Description, "error closing connection to", s.host, s.port, "with message", err)
		s.Ok = false
		return
	}

	log.Println(s.Description, "returned success.")
	s.Ok = true
}
