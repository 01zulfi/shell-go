package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Println("Error while reading command: ", err)
			return
		}

		out := commandNotFoundMesssage(input[:len(input)-1])
		fmt.Fprint(os.Stdout, out)

		fmt.Println()
	}
}

func commandNotFoundMesssage(cmd string) string {
	return fmt.Sprintf("%s: command not found", cmd)
}
