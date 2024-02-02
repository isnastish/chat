package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	Id         string
	Conn       net.Conn
	Addr       net.Addr
	UniqueName string

	IsConnected bool
}

// This is your command registry.
// And then commands.go will be renamed to cli.go and be part of a server package,
// since it's server specific.
// But we can go the other way around and do commands validation on the client side,
// and ONLY send correct commands with their arguments over the network.
// It will help to handle cases like (mkdir <dirname> & cd <dirname> etc.)

type CommandHandler func(args ...string) []byte
type Cli struct {
	commandRegistry map[string]CommandHandler
}

type Server struct {
	Connections map[string]*Peer
	cli         *Cli // maybe not a pointer, didn't have enough time to think about it.
	Mu          sync.Mutex
}

var server *Server

func (cli *Cli) registerCommand(command string, handler CommandHandler) {
	cli.commandRegistry[command] = handler
}

func (cli *Cli) getHandler(command string) (CommandHandler, bool) {
	// I would make it as simple as possible, so it just returns a handle.
	h, exists := cli.commandRegistry[command]
	if exists {
		return h, true
	}
	return nil, false
}

func NewPeer(connection net.Conn, name string) *Peer {
	addr := connection.RemoteAddr()
	return &Peer{
		Addr:        addr,
		Id:          GenSHA256(addr.String()),
		Conn:        connection,
		UniqueName:  name,
		IsConnected: true,
	}
}

func NewServer() *Server {
	s := Server{
		Connections: make(map[string]*Peer),
		cli: &Cli{
			commandRegistry: make(map[string]CommandHandler),
		},
	}

	// Register all the commans that can be executed on a server.
	// We can add custom commands as well, but a handler would be different.
	s.cli.registerCommand(":ls", ls)
	s.cli.registerCommand(":cd", cd)
	s.cli.registerCommand(":cwd", cwd)
	s.cli.registerCommand(":cat", cat)
	s.cli.registerCommand(":mkdir", mkdir)
	s.cli.registerCommand(":rmdir", rmdir)
	s.cli.registerCommand(":rm", rm)
	s.cli.registerCommand(":touch", touch)
	s.cli.registerCommand(":du", du) // disk usage

	return &s
}

func (s *Server) addConnection(conn net.Conn) *Peer {
	log.Println("Added new connection")
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// fmt.Print("Enter unique client name:")

	// // NOTE(alx): This has to be done on the client side,
	// // And the name has to be sent to the server.

	// var uniqueName string
	// // var retriesCount = 3
	// var nameWasSelected = true

	// for {
	// 	fmt.Scanf("%s", &uniqueName)

	// 	// Check in a database whether this name already exists,
	// 	// if not prompt for login.
	// 	for _, peer := range s.Connections {
	// 		if uniqueName == peer.UniqueName {
	// 			fmt.Println("Name already occupied.")
	// 			nameWasSelected = false
	// 			break
	// 		}
	// 	}

	// 	if nameWasSelected {
	// 		break
	// 	}
	// }

	peer := NewPeer(conn, "SomeName")
	s.Connections[peer.Id] = peer
	return peer
}

func (s *Server) removeConnection(connId string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Connections, connId)
}

func (s *Server) processConnection(conn net.Conn) {

	curPeer := s.addConnection(conn) // pointer might change its address as we add more connections.
	log.Println("Connected peer: ", curPeer.Addr.String())

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Failed to close the connection: ", curPeer.Addr.String())
		}
		s.removeConnection(curPeer.Id)
		fmt.Println("Peer disconnected: ", curPeer.Addr.String())
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		data := strings.Split(input.Text(), " ")
		command := strings.ToLower(TrimWhitespaces(data[0]))
		args := data[1:]

		// NOTE: This is not a final version, because if we want to handle edge cases like this:
		// mkdir <dirname> & cd <dirname> we would have to parse the whole string and break it down into tokens.
		if handler, exist := s.cli.getHandler(command); exist {
			log.Println("Invoking :ls command")
			conn.Write(handler(args...))
		} else {
			switch {

			case MatchCommand(command, ":ftp"):
				// NOTE: get file name from the host, returns only bytes for now.
				if len(args) == 0 {
					conn.Write([]byte("File is not specified\n"))
					continue
				}

				f, err := os.Open(args[0])
				if err != nil {
					conn.Write([]byte("File doesn't exist\n"))
				} else {
					defer f.Close()
					io.Copy(conn, f)
				}

			case MatchCommand(command, ":close"):
				// close the connection
				curPeer.IsConnected = false
				return

			case MatchCommand(command, ":echo"):
				// Invoke echo server
				Echo(conn, strings.Join(args, " "), 2*time.Second)

			case MatchCommand(command, ":send"):
				// NOTE(alx): If we want to send file from one client to another,
				// It should come through the server, since we don't support client-to-client communication.
				// Support P2P? The problem with that would be is that client knows nothing about other clients.
				// But on the other hand, it will speed up the process of transferring the files and messages.
				panic("Sending files to different clients is not implemented yet.")

			default:
				// Send messages to all clients.
				for _, peer := range s.Connections {
					if curPeer.Id != peer.Id {
						peer.Conn.Write([]byte(input.Text()))
					}
				}
			}
		}
	}
}

func (s *Server) Run(options *Options) {
	address := options.GetAddress()
	listener, err := net.Listen(options.Network, address)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening: ", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection aborted.")
			continue
		}

		// When client joins the session (and probably session has to be implemented as well.)
		// It should choose the name. The server should store that name in a database or just in memory,
		// use Redis? Try different approaches, with mysql as well.
		go s.processConnection(conn)
	}
}
