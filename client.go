package main

import (
	_ "crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// NOTE(alx): Should it be a part of server instead?

type Session struct {
	// NOTE(alx): Using channel for synchronization is really bad idead!!!
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

// NOTE(alx): Should it be done before we establish the connection?
// But then how do we make sure that the name is unique?
// We would have to make a look up in the database (on the server side).
// TODO(alx): Implement authentication.
func (s *Session) promptForClientName(remote net.Conn) {
	fmt.Print("Enter client name:")

	var (
		uniqueName      string
		totalRetries    = 3
		nameWasSelected = true
	)

	for i := 0; i < totalRetries; i++ {
		fmt.Scanf("%s", &uniqueName)

		remote.Write([]byte("JOIN" + " " + uniqueName))

		if nameWasSelected {
			break
		}
	}

	if !nameWasSelected {
		fmt.Println("Exceeded amount of retries.")
	}
}

func (s *Session) send(dest net.Conn) {
	// promptForClientName()
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

	// Block until either recv or send goroutines invoke send operation on a channel.
	<-session.Done
}
