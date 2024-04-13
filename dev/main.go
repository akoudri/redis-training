package main

import (
	"fmt"

	"akfc.training.com/redis/utils"
)

func main() {

	// Récupérer un article par son ID
	article, err := utils.GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)

	// Récupérer à nouveau le même article (cette fois depuis le cache)
	article, err = utils.GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)
}
