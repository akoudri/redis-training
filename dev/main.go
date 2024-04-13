package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Représente un article de blog
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // pas de mot de passe
		DB:       0,  // utilise la base de données par défaut
	})

	// Délai d'expiration pour le cache (30 secondes dans cet exemple)
	expiration := 30 * time.Second

	// Fonction pour récupérer un article par son ID
	getArticle := func(id int) (*Article, error) {
		// Générer la clé de cache pour l'article
		cacheKey := fmt.Sprintf("article:%d", id)

		// Vérifier si l'article est dans le cache
		val, err := client.Get(context.Background(), cacheKey).Result()
		if err == nil {
			// L'article est dans le cache, le désérialiser et le renvoyer
			var article Article
			if err := json.Unmarshal([]byte(val), &article); err != nil {
				return nil, err
			}
			fmt.Printf("Article %d récupéré depuis le cache\n", id)
			return &article, nil
		}

		// L'article n'est pas dans le cache, le récupérer depuis la base de données
		// (remplacez cette partie par votre logique de récupération de données)
		article := &Article{
			ID:      id,
			Title:   fmt.Sprintf("Article %d", id),
			Content: "Contenu de l'article...",
		}

		// Sérialiser l'article et le stocker dans le cache
		json, err := json.Marshal(article)
		if err != nil {
			return nil, err
		}
		if err := client.Set(context.Background(), cacheKey, json, expiration).Err(); err != nil {
			return nil, err
		}

		fmt.Printf("Article %d récupéré depuis la base de données et mis en cache\n", id)
		return article, nil
	}

	// Récupérer un article par son ID
	article, err := getArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)

	// Récupérer à nouveau le même article (cette fois depuis le cache)
	article, err = getArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)
}
