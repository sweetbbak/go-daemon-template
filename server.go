package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				break
			}
		}

		message := string(buf[0:n])
		fmt.Println("Received message:", message)
		if strings.Contains(message, "exec") {
			c := exec.Command("notify-send", "hello", "from the server")
			c.Run()
		}

		response := "Hello, client! You sent: " + message
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err.Error())
			break
		}
	}
}

func Server(sock string) {
	// sock := "/tmp/unixsock"
	domain := "unix"
	os.Remove(sock) // remove any previous socket file

	var x net.UnixAddr
	x.Name = sock
	x.Net = domain
	l, err := net.ListenUnix("unix", &x)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(sock)
		os.Exit(1)
	}()

	log.Printf("Listening on [%s]...\n", sock)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		fmt.Println("New client connected.")
		go handleConnection(conn)
	}
}
