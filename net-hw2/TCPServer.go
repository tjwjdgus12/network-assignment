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
	"strings"
	"syscall"
)

// to handle ctrl-c
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

	activateSignalHandler()

	serverPort := "22864"
	req_cnt := 0

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	buffer := make([]byte, 1024)

	for {
		// Wait for connection
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

	L1:
		for {
			// Wait for command input
			count, _ := conn.Read(buffer)
			if count == 0 {
				continue
			}
			optionNum := string(buffer[:count])
			fmt.Printf("Command %s\n", optionNum)

			// Process request
			switch optionNum {

			case "1":
				count, _ := conn.Read(buffer)
				conn.Write(bytes.ToUpper(buffer[:count]))

			case "2":
				clientAddr := strings.Split(conn.RemoteAddr().String(), ":")
				response := fmt.Sprintf("clinet IP = %s, port = %s\n", clientAddr[0], clientAddr[1])
				conn.Write([]byte(response))

			case "3":
				response := fmt.Sprintf("requests served = %d\n", req_cnt)
				conn.Write([]byte(response))

			case "5":
				conn.Close()
				req_cnt++
				break L1

			default:
				conn.Write([]byte("Invalid Input!"))
			}

			req_cnt++
		}
		fmt.Printf("\n")
	}
}
