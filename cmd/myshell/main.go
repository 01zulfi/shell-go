package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Input struct {
	command   string
	arguments []string
}

func main() {
	commands := make(map[string]string)
	commands["exit"] = "exit"
	commands["echo"] = "echo"

	for {
		fmt.Fprint(os.Stdout, "$ ")
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading command: ", err)
			return
		}

		text = strings.TrimSpace(text)
		input := Input{
			command:   strings.Split(text, " ")[0],
			arguments: strings.Split(text, " ")[1:],
		}

		switch input.command {
		case commands["exit"]:
			os.Exit(0)
		case commands["echo"]:
			fmt.Fprint(os.Stdout, strings.Join(input.arguments, " "))
		default:
			out := commandNotFoundMesssage(input.command)
			fmt.Fprint(os.Stdout, out)
		}

		fmt.Println()
	}
}

func commandNotFoundMesssage(cmd string) string {
	return fmt.Sprintf("%s: command not found", cmd)
}
