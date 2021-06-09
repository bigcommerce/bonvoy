package main

import (
	"bonvoy/commands"
	"bonvoy/config"
	"os"
)

//func root(args[] string) error {
//	var cmds = commands.All()
//	subcommand := os.Args[1]
//
//	for _, cmd := range cmds {
//		if cmd.Name() == subcommand {
//			var err = cmd.Init(os.Args[2:])
//			if err != nil {
//				panic(err)
//			}
//			return cmd.Run()
//		}
//	}
//	return fmt.Errorf("Unknown subcommand: %s", subcommand)
//}

func main() {
	_ = os.Setenv("DOCKER_API_VERSION", "1.39")
	config.Load()
	commands.Init()
	//if len(os.Args) < 2 {
	//	fmt.Println("USAGE: bonvoy" +
	//		"\n\tversion - Display bonvoy version" +
	//		"\n\tlog-level [service] (level) - Set log level for a sidecar" +
	//		"\n\tcerts-expired - Display expired certificates compared to Consul agent for a sidecar" +
	//		"\n\tlisteners [service] - Display all listeners for sidecar")
	//} else if err := root(os.Args[1:]); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
}