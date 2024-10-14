package main

import (
	"context"
	"log"
	"os"

	"github.com/corka149/rental/cmd"
)

func main() {
	server, err := cmd.NewServer(context.Background(), os.Getenv)

	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	defer server.Close()

	log.Fatalln(server.Run())
}
