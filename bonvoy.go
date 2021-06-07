package main

import (
	"bonvoy/commands"
	"fmt"
	"os"
)

func root(args[] string) error {
	var cmds = commands.All()
	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			var err = cmd.Init(os.Args[2:])
			if err != nil {
				panic(err)
			}
			return cmd.Run()
		}
	}
	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}