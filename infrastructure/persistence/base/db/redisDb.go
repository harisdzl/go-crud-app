package db

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)



func NewRedisDB() (*redis.Client, error) {
	// connection DB
	DbHost := os.Getenv("DB_REDIS_HOST")
	DbPassword := os.Getenv("DB_REDIS_PW")
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