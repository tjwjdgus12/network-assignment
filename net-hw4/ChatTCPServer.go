/**
   ChatTCPServer.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var VERSION = "1.0.0"

type CommandType byte

const (
	CMD_DEFAULT byte = iota
	CMD_LIST         // 1
	CMD_DM           // 2
	CMD_EXIT         // 3
	CMD_VER          // 4
	CMD_RTT          // 5
)

// to handle ctrl-c
func activateSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\ngg~\n\n")
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

	activateSignalHandler()

	serverPort := "22864"
	clientCnt := 0

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n\n", serverPort)

	channel := make(map[string]chan string)

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		buffer := make([]byte, 32)
		count, _ := conn.Read(buffer)
		nickname := string(buffer[:count])

		go func(name string, con net.Conn) {

			fmt.Printf("%s joined from %s. There are %d users connected\n", nickname, conn.RemoteAddr().String(), clientCnt)

			con.Write([]byte("welcome %s to CAU network class chat room at %s. There are %d users connected"))

			for {
				buffer := make([]byte, 1024)
				count, _ := con.Read(buffer)

				command := buffer[0]

				var message string

				switch command {

				case CMD_DEFAULT:
					message = string(buffer[1:count])
					fmt.Println(message)
					con.Write([]byte("ack"))

				case CMD_LIST:
					con.Write([]byte("DUMMY"))

				case CMD_DM:
					message = string(buffer[1:count])
					con.Write([]byte("DUMMY"))

				case CMD_EXIT:
					con.Close()
					fmt.Printf("%s left. There are %d users now\n\n", nickname, clientCnt)
					return

				case CMD_VER:
					con.Write([]byte(VERSION))

				case CMD_RTT:
					con.Write([]byte("DUMMY"))

				default:
					fmt.Print("invalid command\n")
					con.Write([]byte("DUMMY"))
				}
			}
		}(nickname, conn)

		clientCnt++
	}
}
