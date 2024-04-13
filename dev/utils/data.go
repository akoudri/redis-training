package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func ManageCache() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // pas de mot de passe
		DB:       0,  // utilise la base de données par défaut
	})

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

func DisplayArticles() {
	// Récupérer un article par son ID
	article, err := GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)

	// Récupérer à nouveau le même article (cette fois depuis le cache)
	article, err = GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)
}
