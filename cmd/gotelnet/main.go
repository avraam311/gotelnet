package main

import (
	"github.com/avraam311/gotelnet/cmd/app"
	"github.com/avraam311/gotelnet/internal/flags"
	"github.com/avraam311/gotelnet/internal/telnet"
)

func main() {
	flags := flags.New()
	telnet := telnet.New()
	app := app.New(telnet, flags)
	app.Run()
}
