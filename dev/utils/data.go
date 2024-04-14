package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // pas de mot de passe
			DB:       0,  // utilise la base de données par défaut
		})
	}
	return redisClient
}

func StoreString() {
	client := GetRedisClient()

	err := client.Set(context.Background(), "greeting", "Hello", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(context.Background(), "greeting").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Value = ", val)
}

func StoreAndIncrementInt() {
	client := GetRedisClient()
	ctx := context.Background()

	key := "counter"

	increment := func(tx *redis.Tx) error {
		current, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		next := current + 1

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, next, 0)
			return nil
		})
		return err
	}

	_, err := client.SetNX(ctx, key, 0, 0).Result()
	if err != nil {
		panic(err)
	}

	for retries := 10; retries > 0; retries-- {
		err := client.Watch(ctx, increment, key)
		if err == nil {
			break
		}
		if err == redis.TxFailedErr {
			continue
		}
		panic(err)
	}

	val, err := client.Get(ctx, key).Int()
	if err != nil {
		panic(err)
	}
	fmt.Println("Valeur du compteur:", val)
}

func GetData(client *redis.Client, key string) (string, error) {
	// Fonction à tester qui interagit avec Redis
	value, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
