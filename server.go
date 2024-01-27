package main

// NOTE(alx): Sketch.
// if matchCommand("ls") {
// 	ls()
// } else if matchCommand("cd") {
// 	cd()
// }

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var SupportedCommands = map[string]string{
	"ls":    "desc",
	"cd":    "create new directory",
	"close": "close session",
	"get":   "send contents of a file",
	"cat":   "print contents of a file",
	"touch": "<filename>",
	"tree":  "display dir structure",
	"rm":    "remove dir/file(s)",
}

type Client struct {
	Id          string
	Conn        net.Conn
	IsConnected bool
}

type Server struct {
	Connections map[string]*Client
	Mu          sync.Mutex // doesn't require initialization
}

var server *Server

func NewClient(connection net.Conn) *Client {
	return &Client{
		Id:          GenClientId(connection.RemoteAddr().String()),
		Conn:        connection,
		IsConnected: true,
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

func matchCommand(a, b string) bool {
	return bool(a == b)
}

func (s *Server) processCommands() {

}

func (s *Server) processConnection(conn net.Conn) {
	connId := GenClientId(conn.RemoteAddr().String())
	s.addConnection(connId, conn)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		data := strings.Split(input.Text(), " ")
		command := strings.ToLower(TrimWhitespaces(data[0]))
		args := data[1:]
		if matchCommand(command, "ls") {
			cmd := exec.Command("ls.exe", args...)
			log.Printf("Command received: %v\n", data)
			cmdOut, err := cmd.Output()
			CheckError(err)
			conn.Write(cmdOut)
		} else if matchCommand(command, "cd") {
			// TODO(alx): Check whether it's a single directory.
			newDirName := args[0]
			err := os.Mkdir(newDirName, 0755)
			CheckError(err)
			// Make sure the directory has been created.
			out, err := lsDir()
			CheckError(err)
			fmt.Println(out)
		} else if matchCommand(command, "close") {
			conn.Close()
			s.Mu.Lock()
			s.Connections[connId].IsConnected = false
			s.Mu.Unlock()
			break
		} else if matchCommand(command, "get") {
			if len(args) != 1 {
				log.Fatal("Only one file can be sent over the network.")
			}
			fileName := args[0]
			log.Println("Sending file", fileName)

			if DoesFileExist(fileName) {
				f, err := os.Open(fileName)
				CheckError(err)
				defer f.Close()
				bytesSent, err := io.Copy(conn, f)
				CheckError(err)
				fmt.Printf("Bytes sent: %d\n", bytesSent)
			} else {
				fmt.Printf("File [%s] doesn't exist.\n", fileName)
			}
		} else if matchCommand(command, "touch") { // create new file
			if len(args) != 1 {
				log.Fatal("Only one file can be sent over the network.")
			}
			newFileName := TrimWhitespaces(args[0]) // TODO(alx): Introduce file path validation.
			_, err := os.Create(newFileName)
			CheckError(err)
		} else if matchCommand(command, "tree") { // display dir structure on the server side.
			cmd := exec.Command("tree")
			out, err := cmd.Output()
			CheckError(err)
			_, err = conn.Write(out)
			CheckError(err)
		} else if matchCommand(command, "rm") {
			log.Fatal("not implemented yet.")
		} else {
			Echo(conn, input.Text(), 2*time.Second)
		}
	}
	fmt.Println("Client disconnected: ", conn.RemoteAddr().String()) // Use sha?
	defer func() {
		if s.Connections[connId].IsConnected {
			conn.Close()
		}
	}()
	s.removeConnection(connId)
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
