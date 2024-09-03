package grpc

import "fmt"

type ClientConfig struct {
	Port     uint32
	Hostname string
	WorkDir  string
}

func (c ClientConfig) GetServer() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}
