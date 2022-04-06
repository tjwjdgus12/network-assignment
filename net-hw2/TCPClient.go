/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// to remove endline character on any OS
	var endLine string
	if runtime.GOOS == "windows" {
		endLine = "\r\n"
	} else {
		endLine = "\n"
	}

	serverName := "192.168.0.102" //"nsl2.cau.ac.kr"
	serverPort := "22864"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	// Signal check
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Write([]byte("5"))
		conn.Close()
		fmt.Println("\nBye bye~")
		os.Exit(0)
	}()

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

L1:
	for {
		fmt.Printf("Input option: ")
		optionNum, _ := reader.ReadString('\n')
		optionNum = strings.TrimRight(optionNum, endLine)
		requestTime := time.Now()
		conn.Write([]byte(optionNum))

		var elaspedTime time.Duration

		switch optionNum {
		case "1":
			fmt.Printf("Input sentence: ")
			input, _ := reader.ReadString('\n')
			requestTime = time.Now()
			conn.Write([]byte(input))

		case "2":
			break
		case "3":
			break

		case "5":
			conn.Close()
			break L1
		}

		buffer := make([]byte, 1024)
		conn.Read(buffer)
		elaspedTime = time.Since(requestTime)
		fmt.Printf("Reply from server: %s", string(buffer))
		fmt.Printf("RTT = %.3f ms\n\n", float64(elaspedTime.Microseconds())/1000)
	}
}
