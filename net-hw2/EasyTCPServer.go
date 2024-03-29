/**
   TCPServer.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
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
		fmt.Printf("\nBye bye~\n\n")
		os.Exit(0)
	}()
}

// to make 'HH:MM:SS' format
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
	req_cnt := 0 // how many requests are recieved.

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

			case "1": // send text converted to UPPER-case
				count, _ := conn.Read(buffer)
				response = strings.ToUpper(string(buffer[:count]))

			case "2": // send client's IP address and port number
				clientAddr := strings.Split(conn.RemoteAddr().String(), ":")
				response = fmt.Sprintf("clinet IP = %s, port = %s\n", clientAddr[0], clientAddr[1])

			case "3": // send server request count
				response = fmt.Sprintf("requests served = %d\n", req_cnt)

			case "4": // send server running time
				HHMMSS := duration2HHMMSS(time.Since(startTime))
				response = fmt.Sprintf("run time = %s\n", HHMMSS)

			case "5": // close connection
				conn.Close()
				break L1

			default: // exception
				response = "Invalid Input!\n"
			}

			conn.Write([]byte(response))
			req_cnt++
		}
		fmt.Printf("\n")
	}
}
