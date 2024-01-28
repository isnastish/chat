package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
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

func (s *Server) processConnection(conn net.Conn) {
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
		// This is cumbersome!
		data := strings.Split(input.Text(), " ")
		command := strings.ToLower(TrimWhitespaces(data[0]))
		args := data[1:]

		switch {
		case MatchCommand(command, ":ls"):
			bytes := Ls(args)
			conn.Write(bytes)

		case MatchCommand(command, ":cd"):
			bytes := Cd(args[0])
			conn.Write(bytes)

		case MatchCommand(command, ":cwd"):
			bytes := Cwd()
			conn.Write(bytes)

		case MatchCommand(command, ":cat"):
			bytes := Cat(args[0])
			conn.Write(bytes)

		case MatchCommand(command, ":mkdir"):
			bytes := Mkdir(args[0])
			conn.Write(bytes)

		case MatchCommand(command, ":rmdir"):
			bytes := Rmdir(args[0])
			conn.Write(bytes)

		case MatchCommand(command, ":rm"):
			Rm(args)

		case MatchCommand(command, ":touch"):
			filename := TrimWhitespaces(args[0])
			Touch(filename)

		case MatchCommand(command, ":tree"):
			Tree()

		// custom commands:
		case MatchCommand(command, ":get"):
			f := Get(args[0])
			defer f.Close()
			if f != nil {
				io.Copy(conn, f)
			}

		case MatchCommand(command, ":close"):
			peer.IsConnected = false
			return

		default:
			// Let it be echo for now,
			// but here supposed to be more advanced logic.
			Echo(conn, input.Text(), 2*time.Second)
		}
	}
}

func (s *Server) Run(options *Options) {
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
		go s.processConnection(conn)
	}
}
