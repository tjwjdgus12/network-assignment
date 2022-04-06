/**
   TCPServer.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverPort := "22864"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	buffer := make([]byte, 1024)

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		count, _ := conn.Read(buffer)
		conn.Write(bytes.ToUpper(buffer[:count]))
		conn.Close()
	}
}
