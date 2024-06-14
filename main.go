package main

import (
	myhandlers "go-cache-api/handlers" // Rename the package to avoid conflict
	"log"
	"net/http"

	corsHan "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize cache with some data
	// client := cache.GetClient()
	// ctx := context.Background()

	// err := client.Set(ctx, "foo", "bar", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := client.Get(ctx, "foo").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("foo", val)

	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", myhandlers.GetCache).Methods("GET")       // Use myhandlers package
	r.HandleFunc("/cache", myhandlers.SetCache).Methods("POST")            // Use myhandlers package
	r.HandleFunc("/cache/{key}", myhandlers.DeleteCache).Methods("DELETE") // Use myhandlers package
	r.HandleFunc("/ws", myhandlers.WebSocketHandler)                       // Use myhandlers package

	corsHandler := corsHan.CORS(
		corsHan.AllowedOrigins([]string{"*"}),                                       // Allow requests from all origins
		corsHan.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow all HTTP methods
		corsHan.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allow specified headers
	)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(r)))
}
