package initializers

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func InitializeRedisClient() (*redis.Client , error) {
	context := context.Background()
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping(context).Result()
	
	if err != nil {
		return nil, err
	}

	return client, nil
}