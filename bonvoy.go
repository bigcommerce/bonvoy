package main

import (
	"os"

	"bonvoy/commands"
	"bonvoy/config"
)

func main() {
	if os.Getenv("DOCKER_API_VERSION") == "" {
		_ = os.Setenv("DOCKER_API_VERSION", "1.40")
	}
	config.Load()
	commands.NewRegistry().Init()
}
