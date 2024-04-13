package main

import (
	"fmt"

	"akfc.training.com/redis/utils"
)

func main() {

	article, err := utils.GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)

	article, err = utils.GetArticle(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Article: %+v\n", article)
}
