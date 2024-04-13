package utils

import (
	"context"
	"fmt"
	"time"
)

func StoreWithPipeline() {
	client := GetRedisClient()

	// Créer un pipeline
	pipe := client.Pipeline()

	// Ajouter des commandes au pipeline
	incr := pipe.Incr(context.Background(), "counter")
	pipe.Expire(context.Background(), "counter", time.Hour)
	pipe.SAdd(context.Background(), "set", "member1", "member2", "member3")
	sMembers := pipe.SMembers(context.Background(), "set")
	pipe.Set(context.Background(), "key", "value", 0)
	get := pipe.Get(context.Background(), "key")

	// Exécuter le pipeline
	_, err := pipe.Exec(context.Background())
	if err != nil {
		panic(err)
	}

	// Récupérer les résultats des commandes
	fmt.Println("Counter:", incr.Val())
	fmt.Println("Set Members:", sMembers.Val())
	fmt.Println("Key:", get.Val())
}

func StoreWithMSet() {
	client := GetRedisClient()

	// Définir plusieurs clés et valeurs en une seule commande avec MSET
	err := client.MSet(context.Background(), "key1", "value1", "key2", "value2", "key3", "value3").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer plusieurs valeurs en une seule commande avec MGET
	values, err := client.MGet(context.Background(), "key1", "key2", "key3", "nonexistentKey").Result()
	if err != nil {
		panic(err)
	}

	// Afficher les valeurs récupérées
	for i, value := range values {
		if value == nil {
			fmt.Printf("Key%d: (nil)\n", i+1)
		} else {
			fmt.Printf("Key%d: %s\n", i+1, value)
		}
	}
}
