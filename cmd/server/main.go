package main

import (
	"context"
	"fmt"
	"github.com/Amore14rn/faraway/internal/server"
	"github.com/Amore14rn/faraway/pkg/clock"
	"github.com/Amore14rn/faraway/pkg/config"
	"github.com/Amore14rn/faraway/pkg/redis"
	"github.com/Amore14rn/faraway/pkg/utils"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("start server")

	// loading config from file and env
	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// init context to pass config down
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	ctx = context.WithValue(ctx, "clock", clock.SystemClock{})

	cacheInst, err := redis.InitRedisCache(ctx, configInst.CacheHost, configInst.CachePort)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	ctx = context.WithValue(ctx, "cache", cacheInst)

	// seed random generator to randomize order of quotes
	rand.Seed(time.Now().UnixNano())

	// run server
	go func() {
		serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
		if err := server.Run(ctx, serverAddress); err != nil {
			fmt.Println("server error:", err)
		}
	}()

	// Initiate graceful shutdown
	utils.GracefulShutdown(ctx)

	fmt.Println("Server has been gracefully stopped")
}
