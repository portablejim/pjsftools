package main

import (
	"fmt"
	"os"
)

var commands = map[string]func([]string) error{
	"getperms": func(args []string) error { return GetPerms(args, Command{}, Command{}) },
	"setperms": func(args []string) error { return SetPerms(args, Command{}, Command{}, Command{}) },
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage())
		os.Exit(1)
	}
	cmd, ok := commands[os.Args[1]]
	if !ok {
		fmt.Println(usage())
		os.Exit(1)
	}
	cmdErr := cmd(os.Args[2:])
	if cmdErr != nil {
		os.Stderr.WriteString(cmdErr.Error() + "\n")
	}
}

func usage() string {
	s := "Usage: pjsftools [command] [options]\nAvailable commands:\n"
	for k := range commands {
		s += " - " + k + "\n"
	}
	return s
}
