package api

import (
	"bonvoy/envoy"
	"context"
	"fmt"
	"github.com/Devatoria/go-nsenter"
)

func SetLogLevel(config nsenter.Config, level string) bool {
	fmt.Println("curl -X POST "+ envoy.GetHost() + "/logging?level="+level)
	stdout, stderr, err := config.ExecuteContext(context.Background(),"curl", "-X", "POST", envoy.GetHost() + "/logging?level="+level)
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
	fmt.Println(stdout)
	return true
}