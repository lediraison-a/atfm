package server

import (
	"net"
	"net/rpc"
)

type Server struct {
	Inbound *net.TCPListener

	clients map[string]string
}

func (r *Server) Connect(client *string, id *string) error {
	cid := "c" + *client
	*id = cid
	r.clients[*client] = cid
	return nil
}

func NewServer() (*Server, error) {
	port := ":34598"
	addy, err := net.ResolveTCPAddr("tcp", "localhost"+port)
	if err != nil {
		return nil, err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return nil, err
	}
	s := &Server{
		Inbound: inbound,
		clients: map[string]string{},
	}
	fm := NewFileManager()
	rpc.Register(s)
	rpc.Register(fm)
	go rpc.Accept(inbound)

	return s, nil
}
