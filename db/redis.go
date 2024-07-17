package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"password-lock/config"
)

func ConnectToRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to connect to Redis: %s", err))
	} else {
		log.Println("Successfully connected to Redis !")
	}
	return client
}
