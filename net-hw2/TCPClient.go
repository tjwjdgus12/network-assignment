/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	serverName := "nsl2.cau.ac.kr"
	serverPort := "22864"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	fmt.Printf("Input option: ")
	optionNum, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	conn.Write([]byte(optionNum))

	switch optionNum {
	case "1":
		fmt.Printf("Input sentence: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(input))
	case "5":

	}

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Printf("Reply from server: %s", string(buffer))

	conn.Close()
}
