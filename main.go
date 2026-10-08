package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const version = "v7"

var mu sync.Mutex

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	message := os.Getenv("APP_MESSAGE")
	if message == "" {
		message = "couldn't read from env!"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is required")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method=%s path=%s pod=%s", r.Method, r.URL.Path, hostname)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"message": message,
			"pod":     os.Getenv("POD_NAME"),
			"version": version,
		})
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/private", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Ruxsat berildi!",
			"pod":     hostname,
			"version": "v4",
		})
	})

	mux.HandleFunc("/visits", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		n := 0
		if b, err := os.ReadFile("/data/visits"); err == nil {
			n, _ = strconv.Atoi(string(b))
		}
		n++

		if err := os.WriteFile("/data/visits", []byte(strconv.Itoa(n)), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"visits":  n,
			"pod":     os.Getenv("POD_NAME"),
			"version": version,
		})
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server started on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
