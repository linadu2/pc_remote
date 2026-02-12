package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/itchyny/volume-go"
	"github.com/micmonay/keybd_event"
)

// StartServer launches the HTTP server and blocks until it crashes
func StartServer(cfg Config) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Printf("Warning: Failed to init keyboard: %v", err)
	}

	// Linux specific delay for keyboard bonding
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	mux := http.NewServeMux()

	// Middleware for Auth
	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if strings.TrimPrefix(authHeader, "Bearer ") != cfg.Token {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}

	mux.HandleFunc("/action", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Action string `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if req.Action == "play_pause" {
			kb.SetKeys(keybd_event.VK_MEDIA_PLAY_PAUSE)
			if err := kb.Launching(); err != nil {
				log.Printf("Key press error: %v", err)
			} else {
				log.Println("Action: Play/Pause")
			}
		}
		w.WriteHeader(http.StatusOK)
	}))

	mux.HandleFunc("/volume", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req struct {
				Volume int `json:"volume"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			if err := volume.SetVolume(req.Volume); err != nil {
				http.Error(w, "Volume error", http.StatusInternalServerError)
				return
			}
			log.Printf("Volume set to %d", req.Volume)
		} else if r.Method == http.MethodGet {
			vol, _ := volume.GetVolume()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]int{"volume": vol})
		}
	}))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
