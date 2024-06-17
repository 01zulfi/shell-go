package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Input struct {
	command   string
	arguments []string
}

type Shell struct {
	paths    []string
	commands map[string]string
}

func (s *Shell) isBuiltin(command string) bool {
	_, ok := s.commands[command]
	return ok
}

func (s *Shell) inPath(command string) (path string, isInPath bool) {
	pathsMap := s.getPathsMap()
	for path, fileNames := range pathsMap {
		if slices.Contains(fileNames, command) {
			return path + "/" + command, true
		}
	}
	return "", false
}

func (s *Shell) getPathsMap() map[string][]string {
	pathsMap := make(map[string][]string)
	for _, path := range s.paths {
		pathsMap[path] = getAllFileNames(path)
	}
	return pathsMap
}

func newShell() Shell {
	path := os.Getenv("PATH")
	commands := make(map[string]string)
	commands["exit"] = "exit"
	commands["echo"] = "echo"
	commands["type"] = "type"
	commands["pwd"] = "pwd"
	commands["cd"] = "cd"
	shell := Shell{
		commands: commands,
		paths:    strings.Split(path, ":"),
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
			fmt.Println()
		case shell.commands["type"]:
			argument := strings.Join(input.arguments, " ")
			isBuiltin := shell.isBuiltin(argument)
			path, isInPath := shell.inPath(argument)

			var out string
			if isBuiltin {
				out = fmt.Sprintf("%s is a shell builtin", argument)
			} else if isInPath {
				out = fmt.Sprintf("%s is %s", argument, path)
			} else {
				out = fmt.Sprintf("%s: not found", argument)
			}
			fmt.Fprint(os.Stdout, out)
			fmt.Println()
		case shell.commands["pwd"]:
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Fprint(os.Stdout, "error executing pwd")
			}
			fmt.Fprint(os.Stdout, pwd)
			fmt.Println()
		case shell.commands["cd"]:
			argument := strings.Join(input.arguments, " ")
			var err error
			if argument == "~" {
				home, errInside := os.UserHomeDir()
				if errInside != nil {
					err = errInside
				} else {
					err = os.Chdir(home)
				}
			} else {
				err = os.Chdir(argument)
			}
			if err != nil {
				out := fmt.Sprintf("cd: %s: No such file or directory", argument)
				fmt.Fprint(os.Stdout, out)
				fmt.Println()
			}
		default:
			path, isInPath := shell.inPath(input.command)
			if isInPath {
				externalCommand := exec.Command(path, input.arguments...)
				out, err := externalCommand.CombinedOutput()
				if err != nil {
					errorOutput := fmt.Sprintf("error executing %s", input.command)
					fmt.Fprint(os.Stdout, errorOutput)
					break
				}
				fmt.Fprint(os.Stdout, strings.TrimSpace(string(out)))
				fmt.Println()
				break
			}

			out := commandNotFoundMesssage(input.command)
			fmt.Fprint(os.Stdout, out)
			fmt.Println()
		}
	}
}

func commandNotFoundMesssage(cmd string) string {
	return fmt.Sprintf("%s: command not found", cmd)
}

func getAllFileNames(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return make([]string, 1)
	}
	var out []string
	for _, file := range files {
		out = append(out, file.Name())
	}
	return out
}
