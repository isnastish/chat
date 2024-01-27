package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Session struct {
	Done chan struct{}
}

func NewSession() *Session {
	return &Session{
		Done: make(chan struct{}),
	}
}

func (s *Session) recv(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src) // conn net.Conn

	if err != nil {
		fmt.Println("Ignore the error for now.")
	}
	// CheckError(err) // ignore errors for now.
	s.Done <- struct{}{}
}

func Run(options *Options) {
	conn, err := net.Dial(options.Network, options.GetAddress(options.Ports[0]))
	CheckError(err)

	session := NewSession()

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Failed to close the connection.")
		}
	}()

	fmt.Println("Connected to: ", conn.RemoteAddr().String())

	go session.recv(os.Stdout, conn)

	// receive messages from the server on a separate routine. (what's the difference)
	go session.recv(conn, os.Stdin) // try std::err

	<-session.Done
}
