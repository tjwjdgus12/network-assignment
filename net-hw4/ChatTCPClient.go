/**
   ChatTCPClient.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// command list
const (
	CMD_DEFAULT byte = iota
	CMD_LIST         // 1
	CMD_DM           // 2
	CMD_EXIT         // 3
	CMD_VER          // 4
	CMD_RTT          // 5
	CMD_INVALID = 255
)

// to handle ctrl-c
func activateSignalHandler(conn net.Conn) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	conn.Write([]byte{CMD_EXIT})
	conn.Close()
	fmt.Printf("\ngg~\n\n")
	os.Exit(0)
}

// input string -> command byte, message, success (not empty)
func parseInput(input string) (byte, string, bool) {

	input = strings.TrimSpace(input)

	if input == "" {
		return 0, "", false
	}

	message := ""
	command := byte(0)

	if input[0] == '\\' {
		commandStr := ""
		delimIdx := strings.IndexByte(input, ' ')

		if delimIdx == -1 {
			commandStr = input
		} else {
			commandStr = input[:delimIdx]
			message = input[delimIdx+1:]
		}

		switch commandStr {
		case `\list`:
			command = CMD_LIST
		case `\dm`:
			command = CMD_DM
		case `exit`:
			command = CMD_EXIT
		case `\ver`:
			command = CMD_VER
		case `\rtt`:
			command = CMD_RTT
		default:
			command = CMD_INVALID
		}
	} else {
		message = input
	}

	return command, message, true
}

func main() {
	if len(os.Args) != 2 {
		panic("invalid arguments")
	}
	nickname := os.Args[1]

	reader := bufio.NewReader(os.Stdin)

	serverName := "nsl2.cau.ac.kr"
	serverPort := "22864"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	conn.Write([]byte(nickname))

	response := make([]byte, 1024)
	count, _ := conn.Read(response)
	success := response[0] == '1'
	fmt.Println(string(response[1:count]))

	// full room or duplicated nickname
	if !success {
		conn.Close()
		os.Exit(0)
	}

	go activateSignalHandler(conn)

	var rttRequestTime time.Time

	// message(from server) reciever
	go func() {
		data := make([]byte, 1024)
		for {
			count, _ := conn.Read(data)
			if count == 0 {
				continue
			}
			message := string(data[:count])
			if message == "RTT" {
				fmt.Printf("RTT = %.3f ms\n", float64(time.Since(rttRequestTime).Microseconds())/1000)
				continue
			}

			fmt.Println(message)
		}
	}()

	for {
		input, _ := reader.ReadString('\n')
		command, message, ok := parseInput(input)
		if !ok {
			continue
		}
		if command == CMD_INVALID {
			fmt.Println("invalid command")
			continue
		}
		if command == CMD_RTT {
			rttRequestTime = time.Now()
		}

		data := fmt.Sprintf("%c%s", command, message)
		conn.Write([]byte(data))

		if command == CMD_EXIT {
			break
		}

		//fmt.Println()
	}

	fmt.Printf("gg~\n")
}
