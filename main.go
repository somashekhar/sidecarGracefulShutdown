package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownInProgress bool

func main() {
	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for shutdown signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Main application started and ready to serve on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-done
	log.Println("Main container received SIGTERM. Starting graceful shutdown...")
	shutdownInProgress = true

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Simulate cleanup work
	log.Println("Main container: Closing connections...")
	time.Sleep(5 * time.Second)
	log.Println("Main container: Saving state...")
	time.Sleep(5 * time.Second)
	log.Println("Main container: Cleanup complete")

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if shutdownInProgress {
		fmt.Fprint(w, "Application is shutting down...\n")
		return
	}
	fmt.Fprint(w, "Hello from main container!\n")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if shutdownInProgress {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Not healthy, shutting down")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Healthy")
}