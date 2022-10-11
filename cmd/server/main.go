package main

import (
	"log"
	"net"
	"net/http"
)

const (
	HOST = "0.0.0.0"
	PORT = "8000"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "🛑 not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"msg":"🍌"}`))
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🏥 ok"))
	})

	addr := net.JoinHostPort(HOST, PORT)
	log.Println("🍌 running on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
