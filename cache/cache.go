package cache

import "github.com/redis/go-redis/v9"

var client *redis.Client

func GetClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		return client
	}
	return client
}
