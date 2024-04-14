package utils

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestGetData(t *testing.T) {
	// Démarrer le serveur miniredis
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Erreur lors du démarrage de miniredis : %v", err)
	}
	defer s.Close()

	// Configurer le client Redis pour utiliser miniredis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	// Préparer les données de test dans miniredis
	s.Set("key1", "value1")
	s.Set("key2", "value2")

	// Appeler la fonction à tester
	result, err := GetData(client, "key1")
	if err != nil {
		t.Errorf("Erreur inattendue : %v", err)
	}

	// Vérifier le résultat
	expected := "value1"
	if result != expected {
		t.Errorf("Résultat incorrect. Attendu : %s, Obtenu : %s", expected, result)
	}

	// Vérifier les interactions avec Redis
	if !s.Exists("key1") {
		t.Error("La clé 'key1' devrait exister dans Redis")
	}
	if !s.Exists("key2") {
		t.Error("La clé 'key2' devrait exister dans Redis")
	}
}
