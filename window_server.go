//go:build windows
// +build windows

package main

import (
	"fmt"

	"github.com/micmonay/keybd_event"
)

// Windows specific constants
const (
	KeyPlayPause = keybd_event.VK_MEDIA_PLAY_PAUSE
	KeyNext      = keybd_event.VK_MEDIA_NEXT_TRACK
	KeyPrev      = keybd_event.VK_MEDIA_PREV_TRACK
	KeyStop      = keybd_event.VK_MEDIA_STOP
)

func shutdownSystem() error {
	// /s = shutdown, /t 0 = immediately
	//return exec.Command("shutdown", "/s", "/t", "0").Run()
	fmt.Println("Shut down PC")
	return nil
}
