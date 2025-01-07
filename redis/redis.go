package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DSN"),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return client
}

func SetClientConnected(client *redis.Client, key string, value any) error {

	ctx := context.Background()
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("could not set key in Redis: %v", err)
	}
	color.Green("Client %s connected successfully !\n", key)
	return nil
}

func CheckClientConnection(client *redis.Client, key string) bool {
	ctx := context.Background()
	value := client.Get(ctx, key)

	return value.Val() != "" || value.Val() != "0"
}

func CloseClientConnection(client *redis.Client, key string) error {
	if CheckClientConnection(client, key) {
		ctx := context.Background()
		err := client.Set(ctx, key, "0", 0).Err()
		if err != nil {
			return fmt.Errorf("could not close client in Redis: %v", err)
		}
		color.Magenta("Client %s connection closed successfully !\n", key)
	}
	return nil

}
