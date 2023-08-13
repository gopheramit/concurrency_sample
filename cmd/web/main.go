package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Data struct {
	FieldString string `json:"fieldString"`
	Number      int    `json:"number"`
	Process     bool   `json:"process"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode input", http.StatusBadRequest)
		return
	}

	if data.Process {
		sleepTime := 2 + rand.Intn(4)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		data.Number *= 2
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/process", handler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
