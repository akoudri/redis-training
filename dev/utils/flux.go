package utils

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func ExecuteFlux() {
	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	// Nom du flux Redis
	streamName := "mystream"

	// Intervalle d'écriture et de lecture (en millisecondes)
	interval := 1000

	// Contexte pour les opérations Redis
	ctx := context.Background()

	// Lancer la goroutine d'écriture
	go writeToStream(ctx, client, streamName, interval)

	// Lancer la goroutine de lecture
	go readFromStream(ctx, client, streamName, interval)

	// Attendre indéfiniment
	select {}
}

// Fonction pour écrire des données dans le flux Redis
func writeToStream(ctx context.Context, client *redis.Client, streamName string, interval int) {
	//TODO
}

// Fonction pour lire des données depuis le flux Redis
func readFromStream(ctx context.Context, client *redis.Client, streamName string, interval int) {
	// Groupe de consommateurs et nom du consommateur
	// groupName := "mygroup"
	// consumerName := "myconsumer"

	// TODO
}
