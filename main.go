package main

import (
	"context"
	"fmt"

	"go-cache-api/cache"
	"go-cache-api/handlers"
	"log"
	"net/http"

	// "go-cache-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize cache with some data
	client := cache.GetClient()
	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", handlers.GetCache).Methods("GET")
	r.HandleFunc("/cache", handlers.SetCache).Methods("POST")
	r.HandleFunc("/cache/{key}", handlers.DeleteCache).Methods("DELETE")
	r.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
