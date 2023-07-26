package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("start2")
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')
	fmt.Println("hello!,", inputText)
}
