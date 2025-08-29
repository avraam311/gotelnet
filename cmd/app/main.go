package app

import (
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
	a.telnet.ConnectAndServe(a.flags.Host, a.flags.Port, a.flags.Timeout)
}
