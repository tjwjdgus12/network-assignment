/**
   UDPClinet.go
   by Jeong-Hyeon Seo (20172864)
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

	serverName := "nsl2.cau.ac.kr"
	serverPort := "22864"

	pconn, _ := net.ListenPacket("udp", ":")
	server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)

	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n\n", localAddr.Port)

	// Signal check
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pconn.WriteTo([]byte("5"), server_addr)
		pconn.Close()
		fmt.Printf("\nBye bye~\n\n")
		os.Exit(0)
	}()

L1:
	for {

		fmt.Printf("<Menu>\n")
		fmt.Printf("1) convert text to UPPER-case\n")
		fmt.Printf("2) get my IP address and port number\n")
		fmt.Printf("3) get server request count\n")
		fmt.Printf("4) get server running time\n")
		fmt.Printf("5) exit\n")

		fmt.Printf("Input option: ")
		optionNum, _ := reader.ReadString('\n')
		optionNum = strings.TrimRight(optionNum, endLine) // remove endline
		requestTime := time.Now()                         // start time measurement
		pconn.WriteTo([]byte(optionNum), server_addr)

		switch optionNum {

		case "1":
			fmt.Printf("Input sentence: ")
			input, _ := reader.ReadString('\n')
			requestTime = time.Now()
			pconn.WriteTo([]byte(input), server_addr)

		case "5":
			pconn.Close()
			break L1
		}

		buffer := make([]byte, 1024)
		pconn.ReadFrom(buffer)
		elaspedTime := time.Since(requestTime) // end time measurement
		fmt.Printf("Reply from server: %s", string(buffer))
		fmt.Printf("RTT = %.3f ms\n\n", float64(elaspedTime.Microseconds())/1000)
	}

	fmt.Printf("\nBye bye~\n\n")
}
