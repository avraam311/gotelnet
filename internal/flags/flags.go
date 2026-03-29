package flags

import (
	"log"
	"strconv"

	"github.com/spf13/pflag"
)

type Flags struct {
	Host    string
	Port    int
	Timeout int
}

func New() *Flags {
	timeout := pflag.Int("timeout", 10, "timeout")
	pflag.Parse()
	host := pflag.Arg(0)
	portStr := pflag.Arg(1)

	if host == "" || portStr == "" {
		log.Fatal("host or port is empty. Usage: gotelnet [host] [port] [--timeout N]")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid port '%s': %v. Port must be integer.", portStr, err)
	}
	if port < 1 || port > 65535 {
		log.Fatalf("invalid port %d: must be between 1 and 65535.", port)
	}

	flags := Flags{
		Host:    host,
		Port:    port,
		Timeout: *timeout,
	}

	return &flags
}
