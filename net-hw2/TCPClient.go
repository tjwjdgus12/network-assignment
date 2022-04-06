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

	serverName := "10.210.60.90" //"nsl2.cau.ac.kr"
	serverPort := "22864"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

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
		fmt.Printf("\nInput option: ")
		optionNum, _ := reader.ReadString('\n')
		optionNum = strings.TrimRight(optionNum, endLine)
		conn.Write([]byte(optionNum))

		switch optionNum {

		case "1":
			fmt.Printf("Input sentence: ")
			input, _ := reader.ReadString('\n')
			conn.Write([]byte(input))
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("Reply from server: %s", string(buffer))

		case "2":
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("Reply from server: %s", string(buffer))

		case "3":
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("Reply from server: %s", string(buffer))

		case "5":
			conn.Close()
			break L1
		}
	}
}
