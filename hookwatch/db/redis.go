package db

import (
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(uri string) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	Redis = redis.NewClient(opt)
	log.Println("✅ Redis connected")
}
