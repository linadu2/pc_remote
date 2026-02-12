package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// --- LOGGING SETUP ---
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Application starting...")

	// 1. Check for valid existing config
	cfg, valid := LoadConfig()

	// 2. If invalid or missing, Launch GUI
	if !valid {
		log.Println("Config missing or invalid. Launching GUI setup...")
		newCfg, saved := RunGUI()
		if !saved {
			log.Println("Setup cancelled. Exiting.")
			os.Exit(0)
		}
		cfg = newCfg
	}

	// 3. Setup Signal Handling for Ctrl+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// 4. Launch Server in Goroutine
	log.Printf("Starting Server on port %d...", cfg.Port)

	// We need to modify StartServer to accept a context or return the http.Server
	// But since we want to keep it simple, we will just run it and let the main function block on the signal.

	go StartServer(cfg)

	// Block until signal received
	<-stop
	log.Println("Shutting down server...")
	// In a more complex app, you'd call server.Shutdown(ctx) here
	os.Exit(0)
}
