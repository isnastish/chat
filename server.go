package main

import (
	"bufio"
	"fmt"
	_ "io"
	"log"
	"net"
	_ "os"
	_ "os/exec"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Id   string
	Conn net.Conn
}

type Server struct {
	Connections map[string]*Client
	Mu          sync.Mutex // doesn't require initialization
}

var server *Server

func NewClient(connection net.Conn) *Client {
	return &Client{
		Id:   GenClientId(connection.RemoteAddr().String()),
		Conn: connection,
	}
}

func NewServer() *Server {
	return &Server{
		Connections: make(map[string]*Client),
	}
}

func (s *Server) addConnection(id string, conn net.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Connections[id] = NewClient(conn)
}

func (s *Server) removeConnection(connId string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Connections, connId)
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func (s *Server) processConnection(conn net.Conn) {
	connId := GenClientId(conn.RemoteAddr().String())
	s.addConnection(connId, conn)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if strings.Contains(text, "ls") {
			// cmd := exec.Command("ls")
			// fmt.Println(cmd)
			// nbytes, _ := io.Copy(os.Stdout, cmd.Stdin)
			// CheckError(err)
			// fmt.Printf("bytes written: %d\n", nbytes)
			fmt.Println("################# ls command was received .###################")
		} else {
			echo(conn, input.Text(), 2*time.Second)
		}
	}
	fmt.Println("Client disconnected: ", conn.RemoteAddr().String()) // Use sha?
	defer conn.Close()
	s.removeConnection(connId)
	return

	// for {
	// 	message := time.Now().Format("3:04PM\n")
	// 	_, err := io.WriteString(conn, message)
	// 	if err != nil {
	// 		fmt.Println("Client disconnected: ", conn.RemoteAddr().String()) // Use sha?
	// 		defer conn.Close()
	// 		s.removeConnection(connId)
	// 		return
	// 	}
	// 	time.Sleep(1000 * time.Microsecond)
	// }
}

func (s *Server) Run(options *Options) {
	var port string
	if len(options.Ports) != 0 { // design a system to select ports.
		port = options.Ports[0]
	} else {
		fmt.Println("Port is not specified, using the default one: 8080")
		port = "8080"
	}

	address := options.Address + ":" + port
	listener, err := net.Listen(options.Network, address)
	CheckError(err)

	// TODO(alx): Use zap logger.
	log.Println("Listening: ", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection aborted.")
			continue
		}

		fmt.Println("Connected client: ", conn.RemoteAddr().String())
		go s.processConnection(conn)
	}
}
