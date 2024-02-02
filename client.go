package main

/* TODO(alx): Get rid of copy function.
[] Give client some time to retry the connection. (when making a request)
	retriesCount := 0
	while retriesCount <= session.retriesTotal{
		try to establish the connection
		response := makeRequest()
		if response.status == "ConnectionFailed"{

		} else {
			panic(errors.New("Failed to make a request with status: ", response.status))
		}

		if retriesCount <= session.retriesTotal{
			time.Sleep(2 * time.Second)
		}
	}
*/

import (
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
