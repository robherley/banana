package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jaypipes/ghw"
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

		gpuInfo, err := ghw.GPU()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		gpus := make([]string, len(gpuInfo.GraphicsCards))
		for i := range gpuInfo.GraphicsCards {
			gpus[i] = gpuInfo.GraphicsCards[i].String()
		}

		response := map[string]interface{}{
			"gpus": gpus,
			"env":  os.Environ(),
		}

		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🏥 ok"))
	})

	addr := net.JoinHostPort(HOST, PORT)
	log.Println("🍌 running on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
