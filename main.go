package main

import (
	"log"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/routes"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	p, err := base.NewPersistence()

	if err != nil {
		log.Fatal(err)
	}

	router := routes.InitRouter(p)

    router.Run()
}