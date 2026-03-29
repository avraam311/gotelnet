package telnet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type Telnet struct {
	originalTermios *unix.Termios
	rawMode         bool
}

func New() *Telnet {
	return &Telnet{}
}

func (t *Telnet) enableRawMode() error {
	fd := int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil
	}
	t.originalTermios = termios

	newTermios := *termios
	newTermios.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	newTermios.Oflag &^= unix.OPOST
	newTermios.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	newTermios.Cflag |= unix.CS8
	newTermios.Cc[unix.VMIN] = 1
	newTermios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &newTermios); err != nil {
		return fmt.Errorf("tcsetattr: %w", err)
	}
	t.rawMode = true
	return nil
}

func (t *Telnet) disableRawMode() {
	if t.rawMode && t.originalTermios != nil {
		fd := int(os.Stdin.Fd())
		unix.IoctlSetTermios(fd, unix.TCSETS, t.originalTermios)
		t.rawMode = false
	}
}

func (t *Telnet) clearInputBuffer() {
	reader := bufio.NewReader(os.Stdin)
	reader.Discard(1024)
}

func (t *Telnet) ConnectAndServe(host string, port int, timeout int) {
	if err := t.enableRawMode(); err != nil {
		log.Printf("raw mode setup warning: %v", err)
	}

	defer t.disableRawMode()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutdown signal received")
		cancel()
	}()

	backoff := time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Connecting to %s:%d\n", host, port)

			conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), time.Duration(timeout)*time.Second)
			if err != nil {
				log.Printf("connection failed: %v. Retrying in %v", err, backoff)
				time.Sleep(backoff)
				backoff = time.Duration(min(int64(backoff.Seconds()*2), 30)) * time.Second
				continue
			}
			defer conn.Close()

			log.Println("connected to server successfully")
			t.clearInputBuffer()

			readCtx, readCancel := context.WithCancel(ctx)
			writeCtx, writeCancel := context.WithCancel(ctx)

			go func() {
				defer readCancel()
				readFromServer(conn, readCtx)
				writeCancel()
			}()

			writeToServer(conn, writeCtx)
			readCancel()
		}
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func writeToServer(conn net.Conn, ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			input, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				log.Printf("error reading stdin: %v", err)
				return
			}
			if len(input) > 0 {
				conn.Write(input)
			}
		}
	}
}

func readFromServer(conn net.Conn, ctx context.Context) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print(scanner.Text() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("server read error: %v", err)
	}
}
