package main

import (
	"context"
	"fmt"

	"github.com/Amore14rn/faraway/app/internal/client"
	"github.com/Amore14rn/faraway/app/internal/pkg/config"
)

func main() {
	fmt.Println("start client")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// loading config from file and env
	cfg := config.GetConfig()

	// init context to pass config down
	address := fmt.Sprintf("%s:%s", cfg.Client.ClientHost, cfg.Client.ClientPort)

	// run client
	err := client.Run(ctx, address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}
