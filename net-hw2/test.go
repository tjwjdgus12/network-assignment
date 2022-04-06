package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Input: ")
		input, _ := reader.ReadString('\n')
		fmt.Println("Output:", input)
	}
}
