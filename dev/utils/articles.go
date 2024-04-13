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

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetArticle(id int) (*Article, error) {
	cacheKey := fmt.Sprintf("article:%d", id)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	expiration := 30 * time.Second

	val, err := client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var article Article
		if err := json.Unmarshal([]byte(val), &article); err != nil {
			return nil, err
		}
		fmt.Printf("Article %d récupéré depuis le cache\n", id)
		return &article, nil
	}

	connStr := "host=localhost port=5432 user=training password=training dbname=training sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, title, content FROM articles WHERE id = $1"
	row := db.QueryRow(query, id)

	var article Article
	err = row.Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Article %d non trouvé", id)
		}
		return nil, err
	}

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
