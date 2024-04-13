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

// Fonction pour récupérer tous les articles
func GetAllArticles() ([]*Article, error) {
	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // pas de mot de passe
		DB:       0,  // utilise la base de données par défaut
	})

	// Délai d'expiration pour le cache
	expiration := 30 * time.Second

	// Clé de cache pour tous les articles
	cacheKey := "articles:all"

	// Vérifie si les articles sont dans le cache
	val, err := client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// Les articles sont dans le cache, les désérialiser et les renvoyer
		var articles []*Article
		if err := json.Unmarshal([]byte(val), &articles); err != nil {
			return nil, err
		}
		fmt.Println("Articles récupérés depuis le cache")
		return articles, nil
	}

	// Les articles ne sont pas dans le cache, les récupérer depuis la base de données
	// Connexion à la base de données PostgreSQL
	connStr := "host=localhost port=5432 user=training password=training dbname=training sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Requête SQL pour récupérer tous les articles
	query := "SELECT id, title, content FROM articles"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats et les stocker dans un slice d'articles
	var articles []*Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}

	// Vérifier les erreurs après avoir parcouru les résultats
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Sérialiser les articles et les stocker dans le cache
	json, err := json.Marshal(articles)
	if err != nil {
		return nil, err
	}
	if err := client.Set(context.Background(), cacheKey, json, expiration).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Articles récupérés depuis la base de données et mis en cache")
	return articles, nil
}

// Fonction pour récupérer un article par son ID en utilisant un script Lua
func GetArticleWithLua(id int) (*Article, error) {
	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // pas de mot de passe
		DB:       0,  // utilise la base de données par défaut
	})

	// Délai d'expiration pour le cache
	expiration := 30 * time.Second

	// Clé de cache pour l'article
	cacheKey := fmt.Sprintf("article:%d", id)

	// Définir le script Lua
	script := redis.NewScript(`
        local cacheKey = KEYS[1]
        local id = ARGV[1]
        local expiration = ARGV[2]

        local article = redis.call("GET", cacheKey)
        if article then
            return article
        end

        -- L'article n'est pas dans le cache, le récupérer depuis la base de données
        -- (remplacez cette partie par votre logique de récupération de données)
        local articleData = {
            id = id,
            title = "Article " .. id,
            content = "Contenu de l'article " .. id
        }
        article = cjson.encode(articleData)

        redis.call("SET", cacheKey, article, "EX", expiration)

        return article
    `)

	// Exécuter le script Lua
	val, err := script.Run(context.Background(), client, []string{cacheKey}, id, expiration.Seconds()).Result()
	if err != nil {
		return nil, err
	}

	// Désérialiser le résultat du script Lua
	var article Article
	if err := json.Unmarshal([]byte(val.(string)), &article); err != nil {
		return nil, err
	}

	fmt.Printf("Article %d récupéré depuis le cache ou la base de données\n", id)
	return &article, nil
}

func DisplayArticle(id int) {
	// Récupérer un article par son ID
	article, err := GetArticle(id)
	// article, err := GetArticleWithLua(id)
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
