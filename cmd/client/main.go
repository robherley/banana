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
	took := time.Now().Sub(start)

	log.Printf("result: %+v\n", result)
	log.Println("output:", string(result.ModelOutputs))
	log.Println("took:", took.String())
}
