package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// handle incoming connections to our socket / daemon
// this is where we will recieve the messages that we send using the daemon binary
// all you need to do is add a function that "handles" that message when it is recieved
// in this example case it literally we expect a shell command and it is then executed
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
		fmt.Println("bytes:", buf[0:n])

		// handle the incoming message - this could be anything, like daemon specific commands
		handleCommands(message)

		// we write a response, we can change this response based on how our command exits and provide
		// error messages or success messages
		response := "Hello, client! You sent: " + message
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err.Error())
			break
		}
	}
}

func handleCommands(cmd string) string {
	ex := System(cmd)
	if ex == 0 {
		return "cmd successful"
	} else {
		return "cmd not successful"
	}
}

func System(cmd string) int {
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err == nil {
		return 0
	}

	// Figure out the exit code
	if ws, ok := c.ProcessState.Sys().(syscall.WaitStatus); ok {
		if ws.Exited() {
			return ws.ExitStatus()
		}

		if ws.Signaled() {
			return -int(ws.Signal())
		}
	}
	return -1
}

func Server(sock string) {
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
		log.Printf("terminating server, removing socket at [%s]...\n", sock)
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

		go handleConnection(conn)
	}
}
