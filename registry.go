package main

import "strings"

type ArrayFlags []string

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
