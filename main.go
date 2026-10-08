package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "v6"

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

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server started on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
