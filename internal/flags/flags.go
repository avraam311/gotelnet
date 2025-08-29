package flags

import (
	"log"

	"github.com/spf13/pflag"
)

type Flags struct {
	Host    string
	Port    string
	Timeout int
}

func New() *Flags {
	timeout := pflag.Int("timeout", 10, "timeout")
	pflag.Parse()
	host := pflag.Arg(0)
	port := pflag.Arg(1)

	if host == "" || port == "" {
		log.Fatal("host or port is empty")
	}

	flags := Flags{
		Host:    host,
		Port:    port,
		Timeout: *timeout,
	}

	return &flags
}
