package handlers

import (
	"context"
	"encoding/json"
	"go-cache-api/cache"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type PubSubMessage struct {
	Event    string `json:"event"`
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
	TimeLeft int64  `json:"time_left"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	client := cache.GetClient()
	ctx := context.Background()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	pubsub := client.PSubscribe(ctx, "__key*__:*")
	defer pubsub.Close()

	done := make(chan bool)

	// Goroutine to handle messages from the WebSocket client
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				done <- true
				break
			}

			// if string(message) == "get_all_data" {
			// 	log.Print("'hoi dslkfkdnfkjn'")
			// 	data := make(map[string]string)

			// 	keys, err := client.Keys(ctx, "*").Result()
			// 	if err != nil {
			// 		log.Printf("Error retrieving keys from Redis: %v\n", err)
			// 		return
			// 	}

			// 	for _, key := range keys {
			// 		value, err := client.Get(ctx, key).Result()
			// 		if err != nil {
			// 			log.Printf("Error retrieving value for key %s: %v\n", key, err)
			// 			continue
			// 		}
			// 		data[key] = value
			// 	}

			// 	log.Printf("Data sent via WebSocket: %v\n", data)
			// 	if err := conn.WriteJSON(data); err != nil {
			// 		log.Printf("Error writing JSON to WebSocket connection: %v\n", err)
			// 		return
			// 	}
			// } else {
			// 	log.Printf("Unknown command: %s\n", message)
			// 	if err := conn.WriteMessage(websocket.TextMessage, []byte("Unknown command")); err != nil {
			// 		log.Printf("Error writing message to WebSocket connection: %v\n", err)
			// 		break
			// 	}
			// }

			if string(message) == "get_all_data" {
				log.Print("'hoi dslkfkdnfkjn'")
				var data []PubSubMessage

				keys, err := client.Keys(ctx, "*").Result()
				if err != nil {
					log.Printf("Error retrieving keys from Redis: %v\n", err)
					return
				}

				for _, key := range keys {
					value, err := client.Get(ctx, key).Result()
					if err != nil {
						log.Printf("Error retrieving value for key %s: %v\n", key, err)
						continue
					}
					ttl, _ := client.TTL(ctx, key).Result()
					// Convert TTL duration to milliseconds
					timeLeftMilliseconds := ttl.Milliseconds()

					data = append(data, PubSubMessage{
						Key:      key,
						Value:    value,
						TimeLeft: timeLeftMilliseconds,
					})
				}

				log.Printf("Data sent via WebSocket: %v\n", data)
				if err := conn.WriteJSON(data); err != nil {
					log.Printf("Error writing JSON to WebSocket connection: %v\n", err)
					return
				}
			} else {
				log.Printf("Unknown command: %s\n", message)
				if err := conn.WriteMessage(websocket.TextMessage, []byte("Unknown command")); err != nil {
					log.Printf("Error writing message to WebSocket connection: %v\n", err)
					break
				}
			}

		}
	}()

	// Goroutine to handle messages from Redis PubSub
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message from Redis PubSub: %v\n", err)
				done <- true
				break
			}
			log.Printf("Received PubSub message: %s - %s\n", msg.Channel, msg.Payload)

			var data PubSubMessage
			switch msg.Channel {
			case "__keyevent@0__:set":
				value, _ := client.Get(ctx, msg.Payload).Result()
				ttl, _ := client.TTL(ctx, msg.Payload).Result()
				timeLeftMilliseconds := ttl.Milliseconds()

				// if err != nil {
				// 	log.Fatalf("Error getting TTL for key %s: %v", key, err)
				// }
				data = PubSubMessage{
					Event:    "set",
					Key:      msg.Payload,
					Value:    value,
					TimeLeft: timeLeftMilliseconds,
				}
			case "__keyevent@0__:expired":
				data = PubSubMessage{
					Event: "expired",
					Key:   msg.Payload,
				}
			case "__keyevent@0__:del":
				data = PubSubMessage{
					Event: "deleted",
					Key:   msg.Payload,
				}
			default:
				continue
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling JSON: %v\n", err)
				continue
			}
			log.Printf("Sending PubSub message: %s\n", jsonData)
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("Error writing PubSub message to WebSocket connection: %v\n", err)
				done <- true
				break
			}
		}
	}()

	<-done
}
