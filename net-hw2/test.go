package main

import (
	"fmt"
)

func main() {
	optionNum := ""
	fmt.Scanln("%s", &optionNum)
	if optionNum == "a" {
		fmt.Println("good")
	}
}
