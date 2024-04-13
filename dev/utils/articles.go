package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

// Représente un article de blog
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Fonction pour récupérer un article par son ID
func GetArticle(id int) (*Article, error) {
	// Générer la clé de cache pour l'article
	cacheKey := fmt.Sprintf("article:%d", id)

	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // pas de mot de passe
		DB:       0,  // utilise la base de données par défaut
	})

	// Délai d'expiration pour le cache
	expiration := 30 * time.Second

	// Vérifie si l'article est dans le cache
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
	// Connexion à la base de données PostgreSQL
	connStr := "host=localhost port=5432 user=training password=training dbname=training sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Requête SQL pour récupérer l'article par son ID
	query := "SELECT id, title, content FROM articles WHERE id = $1"
	row := db.QueryRow(query, id)

	// Scan du résultat dans un objet Article
	var article Article
	err = row.Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Article %d non trouvé", id)
		}
		return nil, err
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
	return &article, nil
}
