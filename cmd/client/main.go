package main

import (
	"log"
	"os"
	"time"

	"github.com/bananaml/banana-go"
)

func main() {
	model := os.Getenv("MODEL_KEY")
	if model == "" {
		log.Fatalln("missing MODEL_KEY env var")
	}

	key := os.Getenv("API_KEY")
	if key == "" {
		log.Fatalln("missing API_KEY env var")
	}

	start := time.Now()
	result, err := banana.Run(key, model, []byte("{}"))
	if err != nil {
		log.Fatalln(err)
	}
	took := time.Since(start)

	log.Println("ID:", result.ID)
	log.Println("Created:", result.Created)
	log.Println("Message:", result.Message)
	log.Println("Took:", took.String())
	log.Println("Output:", string(result.ModelOutputs))
}
