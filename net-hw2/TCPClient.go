/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	var endLine string
	if runtime.GOOS == "windows" {
		endLine = "\r\n"
	} else {
		endLine = "\n"
	}

	serverName := "10.210.60.90" //"nsl2.cau.ac.kr"
	serverPort := "22864"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	fmt.Printf("Input option: ")
	optionNum, _ := reader.ReadString('\n')
	optionNum = strings.TrimRight(optionNum, endLine)
	conn.Write([]byte(optionNum))

	for {
		switch optionNum {
		case "1":
			fmt.Printf("Input sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			conn.Write([]byte(input))
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("Reply from server: %s", string(buffer))
		case "5":
			conn.Close()
			break
		}
	}
}
