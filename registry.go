package main

import "strings"

type ArrayFlags []string

type CmdHandler func(args []string) []byte

type CmdRegistry struct {
	commands map[string]CmdHandler
}

func NewCmdRegistry() *CmdRegistry {
	return &CmdRegistry{
		commands: make(map[string]CmdHandler),
	}
}

func (r *CmdRegistry) RegisterCommand(cmdName string, handler CmdHandler) {
	r.commands[cmdName] = handler
}

func (r *CmdRegistry) ExecuteCommand(cmdName string, args []string) []byte {
	handler, handlerExists := r.commands[cmdName]
	if !handlerExists {
		return []byte("Unknown command\n\n")
	}

	return handler(args)
}

// NOTE(alx): This should only contain the data that we want to keep.
type Options struct {
	Ports   ArrayFlags
	Address string
	Network string
}

func (flags *ArrayFlags) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}

func (flags *ArrayFlags) String() string {
	return strings.Join(*flags, " ")
}

func (o *Options) GetAddress(port string) string {
	return o.Address + ":" + port
}
