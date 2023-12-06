package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

// give function a path to a unix socket to send a message to
func Client(sock string, message string) error {
	conn, err := net.Dial("unix", sock)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Sending message:", message)
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Println("Error sending message:", err.Error())
		return err
	}

	// Wait for response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error receiving response:", err.Error())
		return err
	}
	log.Println("Received response:", string(buf[0:n]))
	return nil
}
