package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Salom! Kubernetes ichidagi Go ilovasidan xabar.")
	})

	fmt.Println("Server 8080-portda ishlamoqda...")
	http.ListenAndServe(":8080", nil)
}
