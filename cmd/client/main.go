package main

import (
	"context"
	"fmt"
	"github.com/Amore14rn/faraway/internal/client"
	"github.com/Amore14rn/faraway/pkg/config"
	"github.com/Amore14rn/faraway/pkg/utils"
)

func main() {
	fmt.Println("start client")

	// loading config from file and env
	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// init context to pass config down
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Run client in a goroutine
	go func() {
		clientAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
		if err := client.Run(ctx, clientAddress); err != nil {
			fmt.Println("client error:", err)
		}
	}()

	// Initiate graceful shutdown
	utils.GracefulShutdown(ctx)

	fmt.Println("Client has been gracefully stopped")
}
