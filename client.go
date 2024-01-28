package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type RemoteConn struct {
	Conn net.Conn
	Addr string
}

type Session struct {
	Done   chan struct{}
	Remote RemoteConn
}

func NewSession(remote net.Conn) *Session {
	return &Session{
		Done:   make(chan struct{}),
		Remote: RemoteConn{remote, remote.RemoteAddr().String()},
	}
}

func (s *Session) recv(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)

	if err != nil {
		fmt.Println("Ignore the error for now.")
	}
	// CheckError(err) // ignore errors for now.
	s.Done <- struct{}{}
}

func Run(options *Options) {
	conn, err := net.Dial(options.Network, options.GetAddress(options.Ports[0]))
	CheckError(err)

	session := NewSession(conn)

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Failed to close the connection.")
		}
	}()

	fmt.Println("Connected to: ", session.Remote.Addr)

	go session.recv(os.Stdout, conn)
	go session.recv(conn, os.Stdin)

	<-session.Done
}
