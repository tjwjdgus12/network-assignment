/**
   TCPServer.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// to check ctrl-c
func activateSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nBye bye~")
		os.Exit(0)
	}()
}

func main() {
	serverPort := "22864"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	activateSignalHandler()

	buffer := make([]byte, 1024)

	for {
		// Wait for connection
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

	L1:
		for {
			// Wait for command input
			count, _ := conn.Read(buffer)
			commandNum, _ := strconv.Atoi(string(buffer[:count]))

			// Process request
			switch commandNum {
			case 1:
				count, _ := conn.Read(buffer)
				conn.Write(bytes.ToUpper(buffer[:count]))
			case 5:
				break L1
			}
		}

		conn.Close()
	}
}
