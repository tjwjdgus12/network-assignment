/**
   P2POmokServer.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Client struct {
	name       string
	address    string
	connection net.Conn
}

const CMD_EXIT byte = 3

func activateSignalHandler(player *[2]Client) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	for i := 0; i <= 1; i++ {
		if player[i].connection != nil {
			player[i].connection.Close()
		}
	}
	fmt.Println("Bye~")
	os.Exit(0)
}

func waitPlayer(listener *net.Listener, cnt *int, player *[2]Client, myNum int, match chan int) {
	opNum := (myNum + 1) % 2

	for {
		conn, _ := (*listener).Accept()

		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)
		data := string(buffer[:count])
		dataList := strings.Split(data, " ")

		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
		player[myNum] = Client{dataList[0], ip + ":" + dataList[1], conn}

		*cnt += 1

		fmt.Printf("\n%s joined from %s. UDP port %s.\n", dataList[0], conn.RemoteAddr().String(), dataList[1])
		if *cnt == 1 {
			fmt.Printf("1 user connected, waiting for another\n")
			conn.Write([]byte("waiting for an opponent."))
		}
		if *cnt == 2 {
			fmt.Printf("2 users connected, notifying %s and %s.\n", player[myNum].name, player[opNum].name)
			msg := fmt.Sprintf("%s is waiting for you (%s).", player[opNum].name, player[opNum].address)
			conn.Write([]byte(msg))
			msg = fmt.Sprintf("%s joined (%s).", player[myNum].name, player[myNum].address)
			player[opNum].connection.Write([]byte(msg))
			match <- 1
		}

		// wait client's sign
		conn.Read(buffer)

		conn.Close()
		*cnt -= 1

		if buffer[0] == CMD_EXIT {
			fmt.Printf("%s disconnected.\n", player[myNum].name)
		}
	}
}

func main() {
	serverPort := "52864"
	listener, _ := net.Listen("tcp", ":"+serverPort)

	clientCnt := 0
	var player [2]Client

	go activateSignalHandler(&player)

	match := make(chan int)

	go waitPlayer(&listener, &clientCnt, &player, 0, match)
	go waitPlayer(&listener, &clientCnt, &player, 1, match)

	for {
		<-match

		time.Sleep(time.Millisecond * 30)

		fmt.Printf("%s and %s disconnected.\n", player[0].name, player[1].name)

		for selfNum := 0; selfNum <= 1; selfNum++ {
			opponentNum := (selfNum + 1) % 2
			dataList := []string{`\play`, player[opponentNum].name, player[opponentNum].address, strconv.Itoa(selfNum + 1)}
			data := strings.Join(dataList, " ")
			player[selfNum].connection.Write([]byte(data))
		}

		time.Sleep(time.Millisecond * 30)
	}
}
