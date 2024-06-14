package models

type CacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Time  string `json:"time"` // New field for time unit
}
