package main

// TODO(alx): Read up on how to test client-server applications.
// NOTE(alx): Client can be executed in one docker container and the server in another.

import (
	"os/exec"
	"testing"
)

func runTestingServer() {
	options := Options{}

	s := NewServer()
	s.Run(&options)
}

func TestServerCommands(t *testing.T) {
	// Boot up the server. Maybe use a separte go routine instead of running a process.
	// Connect to the server with the client.
	// Send commands over the channel.
	// dirName, err := os.MkdirTemp("", "tmp_test")
	// assert.Equal(t, err, nil)
	// defer os.RemoveAll(dirName)

	go runTestingServer()

	// move to that directory.
	// os.Chdir(dirName)

	exec.Command("go", "build")
}
