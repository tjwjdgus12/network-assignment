/**
   P2POmokClient.go
   by Jeong-Hyeon Seo (20172864)
**/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// command list
const (
	CMD_DEFAULT byte = iota
	CMD_PLACE        // 1
	CMD_GG           // 2
	CMD_EXIT         // 3
	CMD_TIMEOUT      // 4
	CMD_INVALID = 255
)

const (
	Row = 10
	Col = 10
)

type Board [][]int

func printBoard(b Board) {
	fmt.Print("   ")
	for j := 0; j < Col; j++ {
		fmt.Printf("%2d", j)
	}

	fmt.Println()
	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()

	for i := 0; i < Row; i++ {
		fmt.Printf("%d |", i)
		for j := 0; j < Col; j++ {
			c := b[i][j]
			if c == 0 {
				fmt.Print(" +")
			} else if c == 1 {
				fmt.Print(" 0")
			} else if c == 2 {
				fmt.Print(" @")
			} else {
				fmt.Print(" |")
			}
		}

		fmt.Println(" |")
	}

	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()
}

func checkWin(b Board, x, y int) int {
	lastStone := b[x][y]
	startX, startY, endX, endY := x, y, x, y

	// Check X
	for startX-1 >= 0 && b[startX-1][y] == lastStone {
		startX--
	}
	for endX+1 < Row && b[endX+1][y] == lastStone {
		endX++
	}

	if endX-startX+1 >= 5 {
		return lastStone
	}

	// Check Y
	startX, startY, endX, endY = x, y, x, y
	for startY-1 >= 0 && b[x][startY-1] == lastStone {
		startY--
	}
	for endY+1 < Row && b[x][endY+1] == lastStone {
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 1
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && startY-1 >= 0 && b[startX-1][startY-1] == lastStone {
		startX--
		startY--
	}
	for endX+1 < Row && endY+1 < Col && b[endX+1][endY+1] == lastStone {
		endX++
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 2
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && endY+1 < Col && b[startX-1][endY+1] == lastStone {
		startX--
		endY++
	}
	for endX+1 < Row && startY-1 >= 0 && b[endX+1][startY-1] == lastStone {
		endX++
		startY--
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	return 0
}

type NetworkStatus struct {
	network string
	conn    net.Conn
	pconn   net.PacketConn
	addr    *net.UDPAddr
}

// to handle ctrl-c
func activateSignalHandler(status *NetworkStatus) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	if status.network == "tcp" {
		status.conn.Write([]byte{CMD_EXIT})
		status.conn.Close()
	} else {
		status.pconn.WriteTo([]byte{CMD_EXIT}, status.addr)
	}
	fmt.Println("Bye~")
	os.Exit(0)
}

// input string -> command byte, message, success (not empty)
func parseInput(input string) (byte, string, bool) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, "", false
	}

	message := ""
	command := CMD_DEFAULT

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
		case `\\`:
			command = CMD_PLACE
		case `\gg`:
			command = CMD_GG
		case `\exit`:
			command = CMD_EXIT
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

	serverName := "nsl2.cau.ac.kr"
	serverPort := "22864"

	pconn, _ := net.ListenPacket("udp", ":")
	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	fmt.Printf("welcome %s to p2p-omok server at %s.\n", nickname, conn.RemoteAddr().String())

	networkStatus := NetworkStatus{"tcp", conn, nil, nil}

	go activateSignalHandler(&networkStatus)

	localAddr := pconn.LocalAddr().(*net.UDPAddr)

	dataList := []string{nickname, strconv.Itoa(localAddr.Port)}
	data := strings.Join(dataList, " ")
	conn.Write([]byte(data))

	buffer := make([]byte, 1024)
	count := 0

	for {
		count, _ = conn.Read(buffer)
		if count == 0 {
			continue
		}
		msg := string(buffer[:count])
		if msg[:5] == `\play` {
			break
		}
		fmt.Println(msg)
	}

	data = string(buffer[:count])
	dataList = strings.Split(data, " ")

	opponentName := dataList[1]
	opponentAddrStr := dataList[2]
	myNum, _ := strconv.Atoi(dataList[3])

	opponentAddr, _ := net.ResolveUDPAddr("udp", opponentAddrStr)
	networkStatus = NetworkStatus{"udp", nil, pconn, opponentAddr}
	opponentNum := (myNum % 2) + 1

	conn.Write([]byte("o"))
	conn.Close()

	reader := bufio.NewReader(os.Stdin)

	board := Board{}
	stoneCnt := 0

	isFinish := false
	myTurn := (myNum == 1)

	if myTurn {
		fmt.Printf("you play first.\n")
	} else {
		fmt.Printf("%s plays first.\n", opponentName)
	}

	for i := 0; i < Row; i++ {
		var tempRow []int
		for j := 0; j < Col; j++ {
			tempRow = append(tempRow, 0)
		}
		board = append(board, tempRow)
	}

	printBoard(board)

	// Reciever
	go func() {
		buffer := make([]byte, 1024)
		for {
			count, _, _ := pconn.ReadFrom(buffer)
			command := buffer[0]
			message := string(buffer[1:count])

			switch command {
			case CMD_DEFAULT:
				fmt.Printf("%s> %s\n", opponentName, message)

			case CMD_PLACE:
				xy := strings.Split(message, " ")
				x, _ := strconv.Atoi(xy[0])
				y, _ := strconv.Atoi(xy[1])

				board[x][y] = opponentNum

				printBoard(board)

				win := checkWin(board, x, y)
				if win != 0 {
					fmt.Println("you lose.")
					isFinish = true
				}

				stoneCnt += 1
				if stoneCnt == Row*Col {
					fmt.Println("draw.")
					isFinish = true
				}

				myTurn = !myTurn

				go func(prevStoneCnt int) {
					<-time.After(time.Second * 10)
					if prevStoneCnt == stoneCnt {
						fmt.Println("time out.")
						fmt.Println("you lose.")
						isFinish = true
						pconn.WriteTo([]byte{CMD_TIMEOUT}, opponentAddr)
					}
				}(stoneCnt)

			case CMD_GG:
				fmt.Printf("%s has given up.\n", opponentName)
				fmt.Println("you win.")
				isFinish = true

			case CMD_EXIT:
				fmt.Printf("%s has exitted.\n", opponentName)
				if !isFinish {
					fmt.Println("you win.")
					isFinish = true
				}

			case CMD_TIMEOUT:
				fmt.Printf("%s has timed out.\n", opponentName)
				fmt.Println("you win.")
				isFinish = true
			}
		}
	}()

	// Sender
	for {
		input, _ := reader.ReadString('\n')
		command, message, ok := parseInput(input)
		if !ok {
			continue
		}
		switch command {

		case CMD_PLACE:
			if isFinish {
				continue
			}
			if !myTurn {
				fmt.Println("not your turn")
				continue
			}
			xy := strings.Split(message, " ")
			if len(xy) != 2 {
				fmt.Println("invalid command")
				continue
			}
			x, errX := strconv.Atoi(xy[0])
			y, errY := strconv.Atoi(xy[1])
			if !(errX == nil && errY == nil) {
				fmt.Println("invalid command")
				continue
			}
			if x < 0 || y < 0 || x >= Row || y >= Col || board[x][y] != 0 {
				fmt.Println("invalid move")
				continue
			}

			board[x][y] = myNum

			printBoard(board)

			win := checkWin(board, x, y)
			if win != 0 {
				fmt.Println("you win.")
				isFinish = true
			}

			stoneCnt += 1
			if stoneCnt == Row*Col {
				fmt.Println("draw.")
				isFinish = true
			}

			myTurn = !myTurn

		case CMD_GG:
			if isFinish {
				continue
			}
			fmt.Println("you have given up.")
			fmt.Println("you lose.")
			isFinish = true

		case CMD_EXIT:
			if !isFinish {
				fmt.Println("you have given up.")
				fmt.Println("you lose.")
				isFinish = true
			}
			pconn.WriteTo([]byte{CMD_EXIT}, opponentAddr)
			fmt.Println("Bye~")
			return

		case CMD_INVALID:
			fmt.Println("invalid command")
			continue
		}

		data := fmt.Sprintf("%c%s", command, message)
		pconn.WriteTo([]byte(data), opponentAddr)
	}
}
