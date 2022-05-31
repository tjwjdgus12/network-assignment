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
	"sync"
	"syscall"
	"time"
)

type Client struct {
	name       string
	connection net.Conn
}

func activateSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\ngg~\n\n")
		os.Exit(0)
	}()
}

func waitPlayer(listener *net.Listener, wg *sync.WaitGroup, player *Client) {
	for {
		conn, _ := (*listener).Accept()

		buffer := make([]byte, 64)
		count, _ := conn.Read(buffer)
		name := string(buffer[:count])

		*player = Client{name, conn}

		wg.Done()

		// wait client's sign
		conn.Read(buffer)
		conn.Close()
		wg.Add(1)
	}
}

func main() {
	serverPort := "22864"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n\n", serverPort)

	var wg sync.WaitGroup
	wg.Add(2)

	var player [2]Client

	go waitPlayer(&listener, &wg, &player[0])
	go waitPlayer(&listener, &wg, &player[1])

	for {
		wg.Wait()

		fmt.Printf("%s vs %s match!\n", player[0].name, player[1].name)

		for selfNum := 0; selfNum <= 1; selfNum++ {
			opponentNum := (selfNum + 1) % 2
			dataList := []string{player[opponentNum].name, player[opponentNum].connection.LocalAddr().String(), strconv.Itoa(selfNum + 1)}
			data := strings.Join(dataList, " ")
			player[selfNum].connection.Write([]byte(data))
		}

		time.Sleep(time.Millisecond)
	}
}
