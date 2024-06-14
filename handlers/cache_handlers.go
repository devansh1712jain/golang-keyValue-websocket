package handlers

import (
	"context"
	"encoding/json"
	"go-cache-api/cache"
	"go-cache-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	client := cache.GetClient() // Get the Redis client instance
	ctx := context.Background()

	value, err := client.Get(ctx, key).Result()
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

func SetCache(w http.ResponseWriter, r *http.Request) {
	var req models.CacheRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := cache.GetClient()
	ctx := context.Background()
	timeValue, _ := strconv.ParseInt(req.Time, 10, 64)

	var expiration time.Duration = time.Duration(timeValue) * time.Second

	err := client.Set(ctx, req.Key, req.Value, expiration).Err()
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	client := cache.GetClient()
	ctx := context.Background()

	deleted, err := client.Del(ctx, key).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
