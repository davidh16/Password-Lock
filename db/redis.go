package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"password-lock/config"
)

func ConnectToRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %s", err))
	}
	fmt.Println("Successfully connected to Redis !")
	return client
}
