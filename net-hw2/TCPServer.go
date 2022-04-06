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
	"time"
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

func duration2HHMMSS(duration time.Duration) string {
	HH := int64(duration.Hours()) % 100
	MM := int64(duration.Minutes()) % 60
	SS := int64(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", HH, MM, SS)
}

func main() {
	startTime := time.Now()

	activateSignalHandler()

	serverPort := "22864"
	req_cnt := 0

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n\n", serverPort)

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
			var response string

			switch optionNum {

			case "1":
				count, _ := conn.Read(buffer)
				conn.Write(bytes.ToUpper(buffer[:count]))

			case "2":
				clientAddr := strings.Split(conn.RemoteAddr().String(), ":")
				response = fmt.Sprintf("clinet IP = %s, port = %s\n", clientAddr[0], clientAddr[1])

			case "3":
				response = fmt.Sprintf("requests served = %d\n", req_cnt)

			case "4":
				HHMMSS := duration2HHMMSS(time.Since(startTime))
				response = fmt.Sprintf("run time = %s\n", HHMMSS)

			case "5":
				conn.Close()
				req_cnt++
				break L1

			default:
				response = "Invalid Input!"
			}

			conn.Write([]byte(response))
			req_cnt++
		}
		fmt.Printf("\n")
	}
}
