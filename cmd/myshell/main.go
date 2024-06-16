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

type Shell struct {
	commands map[string]string
}

func (s *Shell) isBuiltin(command string) bool {
	_, ok := s.commands[command]
	return ok
}

func newShell() Shell {
	commands := make(map[string]string)
	commands["exit"] = "exit"
	commands["echo"] = "echo"
	commands["type"] = "type"
	shell := Shell{
		commands: commands,
	}
	return shell
}

func main() {
	shell := newShell()

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
		case shell.commands["exit"]:
			os.Exit(0)
		case shell.commands["echo"]:
			fmt.Fprint(os.Stdout, strings.Join(input.arguments, " "))
		case shell.commands["type"]:
			argument := strings.Join(input.arguments, " ")
			isBuiltin := shell.isBuiltin(argument)
			var out string
			if isBuiltin {
				out = fmt.Sprintf("%s is a shell builtin", argument)
			} else {
				out = fmt.Sprintf("%s: not found", argument)
			}
			fmt.Fprint(os.Stdout, out)
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
