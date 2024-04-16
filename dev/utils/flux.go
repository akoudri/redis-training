package utils

import (
	"context"
	"fmt"
	"time"

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

func ExecuteFluxWithSubscription() {
	// Créer un client Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	// Nom du flux Redis et du Consumer Group
	streamName := "mystream"
	groupName := "mygroup2"
	consumerName := "myconsumer"

	// Créer le flux s'il n'existe pas déjà
	err := client.XGroupCreate(context.Background(), streamName, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		panic(err)
	}

	// Lancer la goroutine d'écriture
	go writeToStreamWithTimeout(client, streamName)

	// Lancer la goroutine de lecture avec notification
	go readFromStreamWithNotification(client, streamName, groupName, consumerName)

	// Attendre indéfiniment
	select {}
}

// Fonction pour écrire des données dans le flux Redis
func writeToStream(ctx context.Context, client *redis.Client, streamName string, interval int) {
	// Compteur pour générer des valeurs uniques
	counter := 0

	for {
		// Générer une valeur unique
		value := fmt.Sprintf("message %d", counter)

		// Ajouter la valeur au flux Redis
		err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"value": value},
		}).Err()

		if err != nil {
			fmt.Printf("Erreur lors de l'écriture dans le flux : %v\n", err)
		} else {
			fmt.Printf("Valeur écrite dans le flux : %s\n", value)
		}

		// Incrémenter le compteur
		counter++

		// Attendre l'intervalle spécifié avant la prochaine écriture
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func writeToStreamWithTimeout(client *redis.Client, streamName string) {
	// Boucle infinie pour écrire des données dans le flux
	for {
		// Générer une valeur aléatoire
		value := fmt.Sprintf("value %d", time.Now().UnixNano())

		// Ajouter la valeur au flux
		err := client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"data": value},
		}).Err()

		if err != nil {
			fmt.Printf("Erreur lors de l'écriture dans le flux : %v\n", err)
		} else {
			fmt.Printf("Valeur écrite dans le flux : %s\n", value)
		}

		// Attendre 1 seconde avant la prochaine écriture
		time.Sleep(1 * time.Second)
	}
}

// Fonction pour lire des données depuis le flux Redis
func readFromStream(ctx context.Context, client *redis.Client, streamName string, interval int) {
	// Groupe de consommateurs et nom du consommateur
	groupName := "mygroup3"
	consumerName := "myconsumer"

	// Créer le groupe de consommateurs s'il n'existe pas
	err := client.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
	if err != nil && err != redis.Nil {
		fmt.Printf("Erreur lors de la création du groupe de consommateurs : %v\n", err)
		return
	}

	for {
		// Lire les messages du flux Redis
		entries, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    0,
		}).Result()

		if err != nil {
			fmt.Printf("Erreur lors de la lecture depuis le flux : %v\n", err)
		} else {
			for _, entry := range entries[0].Messages {
				fmt.Printf("Valeur lue depuis le flux : %s\n", entry.Values["value"])
			}
		}

		// Attendre l'intervalle spécifié avant la prochaine lecture
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func readFromStreamWithNotification(client *redis.Client, streamName, groupName, consumerName string) {
	// Boucle infinie pour lire les données du flux avec notification
	for {
		// Lire les données du flux avec notification
		entries, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    0, // Bloquer indéfiniment jusqu'à ce que de nouvelles données soient disponibles
		}).Result()

		if err != nil {
			fmt.Printf("Erreur lors de la lecture du flux : %v\n", err)
		} else if len(entries) > 0 {
			// Récupérer la première entrée lue
			entry := entries[0].Messages[0]
			data, ok := entry.Values["data"].(string)
			if !ok {
				fmt.Printf("La valeur 'data' est nil ou n'est pas une chaîne de caractères\n")
				continue
			}
			fmt.Printf("Valeur lue depuis le flux : %s\n", data)

			// Acquitter (acknowledge) le message lu
			err := client.XAck(context.Background(), streamName, groupName, entry.ID).Err()
			if err != nil {
				fmt.Printf("Erreur lors de l'acquittement du message : %v\n", err)
			}
		}
	}
}
