package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/Amore14rn/faraway/app/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedis(ctx context.Context, cfg *config.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()

	return &Redis{
		ctx:    ctx,
		client: rdb,
	}, err
}

func (c *Redis) Add(key int, expiration int64) error {
	return c.client.Set(c.ctx, strconv.Itoa(key), "value", time.Duration(expiration*1e9)*time.Second).Err()
}

func (c *Redis) Get(key int) (bool, error) {
	val, err := c.client.Get(c.ctx, strconv.Itoa(key)).Result()
	return val != "", err
}

func (c *Redis) Delete(key int) {
	c.client.Del(c.ctx, strconv.Itoa(key))
}
