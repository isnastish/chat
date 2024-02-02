package main

import (
	_ "crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Session struct {
	Done chan struct{}
	Conn net.Conn
	Addr string
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Done: make(chan struct{}),
		Conn: conn,
		Addr: conn.RemoteAddr().String(),
	}
}

// Just a temporary
func ignoreErrorForNow(err error) {
	if err != nil {
		fmt.Println("Ignore the error for now.")
	}
}

func (s *Session) recv(src net.Conn) {
	_, err := io.Copy(os.Stdout, src)
	ignoreErrorForNow(err)
	s.Done <- struct{}{}
}

func (s *Session) send(dest net.Conn) {
	_, err := io.Copy(dest, os.Stdin)
	ignoreErrorForNow(err)
	s.Done <- struct{}{}
}

func Run(options *Options) {
	// cert, err := tls.LoadX509KeyPair("generated-cert.pem", "generated-key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // should be part of a server.
	// config := &tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// }

	// conn, err := tls.Dial(options.Network, options.GetAddress(), config)
	conn, err := net.Dial(options.Network, options.GetAddress())
	if err != nil {
		log.Fatal("Connection aborted", err.Error())
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Failed to close the connection.")
		}
	}()

	session := NewSession(conn)

	log.Println("Connected to: ", session.Addr)

	go session.recv(conn)
	go session.send(conn)

	<-session.Done
}
