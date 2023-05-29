package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Amore14rn/faraway/app/internal/pkg/config"
	"github.com/Amore14rn/faraway/app/internal/pkg/redis"
	"github.com/Amore14rn/faraway/app/internal/pkg/timestamp"
	"github.com/Amore14rn/faraway/app/internal/server"
)

func main() {
	fmt.Println("start server")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", timestamp.SystemTime{})

	cfg := config.GetConfig()

	redis, err := redis.NewRedis(ctx, cfg)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	ctx = context.WithValue(ctx, "redis", redis)

	// seed random generator to randomize order of quotes
	rand.Seed(time.Now().UnixNano())

	// run server
	serverAddress := fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)
	err = server.Run(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
