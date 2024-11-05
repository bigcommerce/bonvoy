package main

import (
	"os"

	"bonvoy/commands"
	"bonvoy/config"
)

func main() {
	_ = os.Setenv("DOCKER_API_VERSION", "1.39")
	config.Load()
	commands.NewRegistry().Init()
}
