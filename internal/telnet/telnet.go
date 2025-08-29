package telnet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Telnet struct{}

func New() *Telnet {
	return &Telnet{}
}

func (t *Telnet) ConnectAndServe(host string, port string, timeout int) {
	address := host + ":" + port
	fmt.Println(address)
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatal("error connecting to server - ", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal("error closing connection")
		}
	}()
	log.Print("connected to server successfully")
	go readFromServer(conn)
	writeToServer(conn)
}

func writeToServer(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			log.Fatal("error reading from stdin - ", err)
		}
		if len(input) > 0 {
			_, _ = conn.Write(append(input, '\n'))
		}
		if err == io.EOF {
			break
		}
	}
}

func readFromServer(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading from server - ", err)
	}
}
