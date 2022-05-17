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

func serveClient(myname string, connection map[string]net.Conn) {

	fmt.Printf("%s joined from %s. There are %d users connected\n", myname, connection[myname].RemoteAddr().String(), len(connection))

	for {
		buffer := make([]byte, 1024)
		count, _ := connection[myname].Read(buffer)
		command := buffer[0]

		var message string

		switch command {

		case CMD_DEFAULT:
			message = string(buffer[1:count])
			for target := range connection {
				if target == myname {
					continue
				}
				data := fmt.Sprintf("%s> %s", myname, message)
				connection[target].Write([]byte(data))
			}

		case CMD_LIST:
			data := ""
			for name := range connection {
				data += fmt.Sprintf("%s %s\n", name, connection[name].RemoteAddr().String())
			}
			connection[myname].Write([]byte(data))

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
			data := fmt.Sprintf("from: %s> %s", myname, message)

			if con, ok := connection[target]; ok {
				con.Write([]byte(data))
			}

		case CMD_EXIT:
			connection[myname].Close()
			delete(connection, myname)
			fmt.Printf("%s left. There are %d users now\n", myname, len(connection))
			return

		case CMD_VER:
			connection[myname].Write([]byte(VERSION))

		case CMD_RTT:
			connection[myname].Write([]byte("RTT"))

		default:
			fmt.Print("invalid command\n")
			connection[myname].Write([]byte("invalid command"))
		}

		if strings.Contains(strings.ToUpper(message), "I HATE PROFESSOR") {
			connection[myname].Write([]byte("KILL"))
			connection[myname].Close()
			delete(connection, myname)

			fmt.Printf("[%s is disconnected. There are %d users in the chat room.]\n", myname, len(connection))

			for target := range connection {
				data := fmt.Sprintf("\n[%s is disconnected. There are %d users in the chat room.]", myname, len(connection))
				connection[target].Write([]byte(data))
			}
			return
		}
	}
}

func main() {
	activateSignalHandler()

	serverPort := "22864"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n\n", serverPort)

	connection := make(map[string]net.Conn)

	for {
		conn, _ := listener.Accept()

		buffer := make([]byte, 32)
		count, _ := conn.Read(buffer)
		nickname := string(buffer[:count])

		var response string

		if len(connection) >= 8 { // full room
			response = "0" // fail code
			response += "[chatting room full. cannot connect.]"
			conn.Write([]byte(response))
			continue
		}

		if _, exist := connection[nickname]; exist { // duplicated nickname
			response = "0" // fail code
			response += "[that nickname is already used by another user. cannot connect.]"
			conn.Write([]byte(response))
			continue
		}

		connection[nickname] = conn

		response = "1" // success code
		response += fmt.Sprintf("[welcome %s to CAU network class chat room at %s.]\n", nickname, conn.LocalAddr().String())
		response += fmt.Sprintf("[There are %d users connected.]", len(connection))
		connection[nickname].Write([]byte(response))

		go serveClient(nickname, connection)
	}
}
