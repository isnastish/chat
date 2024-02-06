package main

import "strconv"

type Options struct {
	Port    int
	Address string
	Network string
}

func (o *Options) GetAddress() string {
	ValidatePort(o.Port)
	return o.Address + ":" + strconv.Itoa(o.Port)
}
