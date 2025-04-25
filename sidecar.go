package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	mainAppURL = "http://localhost:8080/health"
	logDir     = "/logs"
)

var shutdownInProgress bool

func main() {
	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Set up signal handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	writeLog("Sidecar started")

	// Start monitoring in goroutine
	go monitorMainApp()

	// Wait for shutdown signal
	<-done
	shutdownInProgress = true
	handleGracefulShutdown()
}

func monitorMainApp() {
	for !shutdownInProgress {
		status := "unhealthy"
		if checkMainApp() {
			status = "healthy"
		}
		writeLog(fmt.Sprintf("Main application is %s", status))
		time.Sleep(5 * time.Second)
	}
}

func checkMainApp() bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(mainAppURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func handleGracefulShutdown() {
	writeLog("Sidecar received SIGTERM. Starting graceful shutdown...")
	writeLog("Sidecar: Starting to monitor main container shutdown")

	// Monitor main app for up to 30 seconds
	for i := 0; i < 30; i++ {
		if !checkMainApp() {
			writeLog("Sidecar: Detected main application is no longer healthy/available")
			break
		}
		writeLog("Sidecar: Main application still running, waiting...")
		time.Sleep(1 * time.Second)
	}

	writeLog("Sidecar: Processing final logs")
	time.Sleep(3 * time.Second) // Simulate final log processing
	writeLog("Sidecar: Shutdown complete, exiting")
}

func writeLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Print to stdout
	fmt.Print(logMessage)

	// Write to file
	logFile := filepath.Join(logDir, "application.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logMessage); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}
}