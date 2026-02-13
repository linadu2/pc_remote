//go:build linux
// +build linux

package main

import (
	"os/exec"

	"github.com/micmonay/keybd_event"
)

// Linux specific constants
const (
	KeyPlayPause = keybd_event.VK_PLAYPAUSE
	KeyNext      = keybd_event.VK_NEXTSONG
	KeyPrev      = keybd_event.VK_PREVIOUSSONG
	KeyStop      = keybd_event.VK_STOPCD
)

func shutdownSystem() error {
	return exec.Command("systemctl", "poweroff").Run()
}
