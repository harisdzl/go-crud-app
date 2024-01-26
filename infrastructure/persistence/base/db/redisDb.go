package db

import (
	"fmt"
	"log"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	"github.com/redis/go-redis/v9"
)



func NewRedisDB() (*redis.Client, error) {
	// connection DB
	DbHost := config.Configuration.GetString("redis.dev.host")
	DbPassword := config.Configuration.GetString("redis.dev.password")
	rdb := redis.NewClient(&redis.Options{
		Addr:	  DbHost,
		Password: DbPassword, // no password set
		DB:		  0,  // use default DB
	})

	if rdb == nil {
		log.Fatal("Failure to connect to Redis")
	} else {
		fmt.Println("Connected to Redis")
	}

	return rdb, nil

}