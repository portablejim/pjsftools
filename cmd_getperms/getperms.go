package cmd_getperms

import (
	"flag"
	"fmt"
	"log"
)

func GetPerms(args []string) {
	log.Printf("GetPerms got args: %v", args)

	fs := flag.NewFlagSet("getperms", flag.ExitOnError)
	orgName := fs.String("org", "", "sf org to use")
	fs.Parse(args)

	if len(*orgName) == 0 {
		fmt.Println("Using default org")
	} else {
		fmt.Println("Using org " + *orgName)
	}
	fmt.Printf("Good morning!\n")
}
