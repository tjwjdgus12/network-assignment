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
	"strings"
	"syscall"
	"time"
)

var VERSION = "1.0.0"

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

func serveClient(name string, con net.Conn, channel map[string]chan string) {

	fmt.Printf("%s joined from %s. There are %d users connected\n", name, con.RemoteAddr().String(), len(channel))

	response := "1" // success code
	response += fmt.Sprintf("[welcome %s to CAU network class chat room at %s.]\n", name, con.LocalAddr().String())
	response += fmt.Sprintf("[There are %d users connected.]", len(channel))
	con.Write([]byte(response))

	// channel recieve && write
	go func() {
		for {
			data := <-channel[name]
			con.Write([]byte(data))
		}
	}()

	for {
		buffer := make([]byte, 1024)
		count, _ := con.Read(buffer)
		command := buffer[0]

		var message string

		switch command {

		case CMD_DEFAULT:
			message = string(buffer[1:count])
			for target := range channel {
				if target == name {
					continue
				}
				data := fmt.Sprintf("%s> %s", name, message)
				channel[target] <- data
			}

		case CMD_LIST:
			data := ""
			for nickname := range channel {
				data += fmt.Sprintf("%s %s\n", nickname, con.RemoteAddr().String())
			}
			channel[name] <- data

		case CMD_DM:
			var target string
			message = string(buffer[1:count])
			delimIdx := strings.IndexByte(message, ' ')
			if delimIdx == -1 {
				target = message
				message = ""
			} else {
				target = message[:delimIdx]
				message = message[delimIdx+1:]
			}
			data := fmt.Sprintf("from: %s> %s", name, message)

			if c, ok := channel[target]; ok {
				c <- data
			}

		case CMD_EXIT:
			con.Close()
			delete(channel, name)
			fmt.Printf("%s left. There are %d users now\n", name, len(channel))
			return

		case CMD_VER:
			con.Write([]byte(VERSION))

		case CMD_RTT:
			con.Write([]byte("RTT"))

		default:
			fmt.Print("invalid command\n")
			con.Write([]byte("invaild command"))
		}

		if strings.Contains(strings.ToUpper(message), "I HATE PROFESSOR") {
			fmt.Printf("[%s is disconnected. There are %d users in the chat room.]\n", name, len(channel))
			for target := range channel {
				if target == name {
					channel[name] <- "KILL"
					continue
				}
				data := fmt.Sprintf("[%s is disconnected. There are %d users in the chat room.]\n", name, len(channel))
				channel[target] <- data
			}
			delete(channel, name)
			con.Close()
			return
		}
	}
}

func main() {
	activateSignalHandler()

	serverPort := "22864"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n\n", serverPort)

	channel := make(map[string]chan string)

	for {
		conn, _ := listener.Accept()

		buffer := make([]byte, 32)
		count, _ := conn.Read(buffer)
		nickname := string(buffer[:count])

		channel[nickname] = make(chan string)

		go serveClient(nickname, conn, channel)
	}
}
