package main

import (
	"fmt"
	"os"
)

var commands = map[string]func([]string){
	"getperms": GetPerms,
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
	cmd(os.Args[2:])
}

func usage() string {
	s := "Usage: pjsftools [command] [options]\nAvailable commands:\n"
	for k := range commands {
		s += " - " + k + "\n"
	}
	return s
}
