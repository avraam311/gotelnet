package app

import (
	"fmt"

	"github.com/avraam311/gotelnet/internal/flags"
	"github.com/avraam311/gotelnet/internal/telnet"
)

type App struct {
	telnet *telnet.Telnet
	flags  *flags.Flags
}

func New(telnet *telnet.Telnet, flags *flags.Flags) *App {
	return &App{
		telnet: telnet,
		flags:  flags,
	}
}

func (a *App) Run() {
	address := fmt.Sprintf("%s:%d", a.flags.Host, a.flags.Port)
	fmt.Println("Telnet client for", address)
	a.telnet.ConnectAndServe(a.flags.Host, a.flags.Port, a.flags.Timeout)
}
