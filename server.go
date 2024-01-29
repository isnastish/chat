package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Peer struct {
	Id          string
	Conn        net.Conn
	Addr        net.Addr
	IsConnected bool
}

type Server struct {
	Connections map[string]*Peer
	Mu          sync.Mutex
}

var server *Server

func NewPeer(connection net.Conn) *Peer {
	addr := connection.RemoteAddr()
	return &Peer{
		Addr:        addr,
		Id:          GenSHA256(addr.String()),
		Conn:        connection,
		IsConnected: true,
	}
}

func NewServer() *Server {
	return &Server{
		Connections: make(map[string]*Peer),
	}
}

func (s *Server) addConnection(conn net.Conn) *Peer {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	peer := NewPeer(conn)
	s.Connections[peer.Id] = peer
	return peer
}

func (s *Server) removeConnection(connId string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Connections, connId)
}

func (s *Server) processConnection(conn net.Conn, cmdReg *CmdRegistry) {
	peer := s.addConnection(conn)
	log.Println("Connected peer: ", peer.Addr.String())

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Failed to close the connection: ", peer.Addr.String())
		}
		s.removeConnection(peer.Id)
		fmt.Println("Peer disconnected: ", peer.Addr.String())
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		data := strings.Split(input.Text(), " ")
		command := strings.ToLower(TrimWhitespaces(data[0]))
		args := data[1:]

		var cmdRes []byte
		if command == ":commands" {
			var tmpRes string
			for cmd := range cmdReg.commands {
				tmpRes += cmd + " "
			}
			cmdRes = []byte("Current available commands:\n" + tmpRes + "\n\n")
		} else {
			cmdRes = cmdReg.ExecuteCommand(command, args)
		}

		conn.Write(cmdRes)
	}
}

func (s *Server) Run(options *Options) {
	cmdReg := NewCmdRegistry()
	cmdReg.RegisterCommand(":ls", Ls)

	var port string
	if len(options.Ports) != 0 {
		port = options.Ports[0]
	} else {
		fmt.Println("Port is not specified, using the default one: 8080")
		port = "8080"
	}

	address := options.GetAddress(port)
	listener, err := net.Listen(options.Network, address)
	CheckError(err)

	log.Println("Listening: ", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection aborted.")
			continue
		}
		go s.processConnection(conn, cmdReg)
	}
}
