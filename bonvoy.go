package main

import (
	"bonvoy/commands"
	"bonvoy/config"
	"os"
)

func main() {
	_ = os.Setenv("DOCKER_API_VERSION", "1.39")
	config.Load()
	commands.NewRegistry().Init()
}