package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

func startServer() {
	cmd := exec.Command("go", "run", "./server/server.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startClient() {
	cmd := exec.Command("go", "run", "./client/client.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting Tunnel Server and Client...")

	// Start the server
	go func() {
		defer wg.Done()
		startServer()
	}()

	// Start the client
	go func() {
		defer wg.Done()
		startClient()
	}()

	wg.Wait()
}
