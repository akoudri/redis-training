package utils

import (
	"context"
	"fmt"
	"time"

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

	expiration := 30 * time.Second

	err := client.Set(context.Background(), "greeting2", "Hello", expiration).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(context.Background(), "greeting").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Value = ", val)
}

func StoreSet() {
	client := GetRedisClient()
	ctx := context.Background()

	// Clé pour l'ensemble
	key := "myset"

	// Ajouter des éléments à l'ensemble
	err := client.SAdd(ctx, key, "item1", "item2", "item3").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les éléments de l'ensemble
	members, err := client.SMembers(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ensemble : ", members)

	// Vérifier si un élément existe dans l'ensemble
	isMember, err := client.SIsMember(ctx, key, "item2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("item2 est un membre : ", isMember)

	// Obtenir le nombre d'éléments dans l'ensemble
	count, err := client.SCard(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Nombre d'éléments : ", count)

	// Supprimer un élément de l'ensemble
	err = client.SRem(ctx, key, "item1").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les éléments de l'ensemble après suppression
	members, err = client.SMembers(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ensemble après suppression : ", members)
}

func StoreList() {
	client := GetRedisClient()
	ctx := context.Background()

	// Clé pour la liste
	key := "mylist"

	// Ajouter des éléments à la liste
	err := client.RPush(ctx, key, "item1", "item2", "item3").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les éléments de la liste
	items, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Liste : ", items)

	// Obtenir la longueur de la liste
	length, err := client.LLen(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Longueur de la liste : ", length)

	// Récupérer un élément spécifique de la liste par son index
	item, err := client.LIndex(ctx, key, 1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Élément à l'index 1 : ", item)

	// Insérer un élément avant un autre élément
	err = client.LInsertBefore(ctx, key, "item2", "item1.5").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les éléments de la liste après l'insertion
	items, err = client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Liste après insertion : ", items)

	// Supprimer un élément de la liste
	err = client.LRem(ctx, key, 1, "item1").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les éléments de la liste après suppression
	items, err = client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Liste après suppression : ", items)
}

func StoreHash() {
	client := GetRedisClient()
	ctx := context.Background()

	// Clé pour le hachage
	key := "myhash"

	// Ajouter des champs et des valeurs au hachage
	err := client.HSet(ctx, key, "field1", "value1", "field2", "value2").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les champs et valeurs du hachage
	fields, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hachage : ", fields)

	// Récupérer la valeur d'un champ spécifique
	value, err := client.HGet(ctx, key, "field1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Valeur du champ 'field1' : ", value)

	// Vérifier si un champ existe dans le hachage
	exists, err := client.HExists(ctx, key, "field2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Le champ 'field2' existe : ", exists)

	// Incrémenter la valeur d'un champ numérique
	err = client.HIncrBy(ctx, key, "field3", 1).Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les champs et valeurs du hachage après l'incrémentation
	fields, err = client.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hachage après incrémentation : ", fields)

	// Supprimer un champ du hachage
	err = client.HDel(ctx, key, "field1").Err()
	if err != nil {
		panic(err)
	}

	// Récupérer tous les champs et valeurs du hachage après suppression
	fields, err = client.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hachage après suppression : ", fields)
}

func StoreAndIncrementInt() {
	client := GetRedisClient()
	ctx := context.Background()

	// Clé pour le compteur
	key := "counter"

	// Fonction pour incrémenter le compteur de manière transactionnelle
	increment := func(tx *redis.Tx) error {
		// Obtenir la valeur actuelle du compteur
		current, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// Incrémenter la valeur
		next := current + 1

		// Exécuter la commande SET
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, next, 0)
			return nil
		})
		return err
	}

	// Instancier le compteur s'il n'existe pas
	_, err := client.SetNX(ctx, key, 0, 0).Result()
	if err != nil {
		panic(err)
	}

	// Incrémenter le compteur de manière transactionnelle
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

	// Obtenir la valeur finale du compteur
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
