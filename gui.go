package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ... [Previous RunGUI code remains exactly the same] ...
// RunGUI launches the configuration window.
func RunGUI() (Config, bool) {
	// ... [Copy previous RunGUI function content here] ...
	// ... [Keep everything the same as before] ...

	// BUT! Since I cannot copy-paste 100 lines again efficiently,
	// I am just providing the REPLACEMENT for the broken functions below.
	// Replace the ENTIRE file content with the previous version
	// BUT replace 'getLocalIP' with the one below.

	// For clarity, here is the full fixed file content:

	a := app.NewWithID("com.pc_control.setup")
	w := a.NewWindow("PC Control Setup")
	w.Resize(fyne.NewSize(400, 350))
	w.SetFixedSize(true)

	// State
	isSaved := false
	var finalConfig Config

	// Data
	portBind := binding.NewInt()
	portBind.Set(8090)
	tokenBind := binding.NewString()
	tokenBind.Set(generateNewToken())

	// UI Components
	tokenEntry := widget.NewEntryWithData(tokenBind)
	tokenEntry.Disable()
	regenBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		tokenBind.Set(generateNewToken())
	})

	copyTokenBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		t, _ := tokenBind.Get()
		w.Clipboard().SetContent(t)
		dialog.ShowInformation("Copied", "Token copied to clipboard!", w)
	})

	portEntry := widget.NewEntryWithData(binding.IntToString(portBind))
	portEntry.Validator = func(s string) error {
		p, err := strconv.Atoi(s)
		if err != nil || p < 1024 || p > 65535 {
			return fmt.Errorf("invalid port")
		}
		return nil
	}

	infoLabel := widget.NewLabel("1. Open Home Assistant\n2. Go to Settings > Devices & Services\n3. Add 'PC Media Controller'\n4. Enter the details below")
	infoLabel.Wrapping = fyne.TextWrapWord
	infoLabel.Alignment = fyne.TextAlignCenter

	// UPDATED IP Logic usage
	ipLabel := widget.NewLabel(fmt.Sprintf("Your IP: %s", getOutboundIP()))
	ipLabel.Alignment = fyne.TextAlignCenter
	ipLabel.TextStyle = fyne.TextStyle{Bold: true}

	saveBtn := widget.NewButton("Save & Run Server", func() {
		p, _ := portBind.Get()
		t, _ := tokenBind.Get()

		finalConfig = Config{Port: p, Token: t}
		if err := SaveConfig(finalConfig); err != nil {
			dialog.ShowError(err, w)
			return
		}

		isSaved = true
		w.Close()
	})
	saveBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewLabelWithStyle("Setup PC Control Server", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		infoLabel,
		ipLabel,
		layout.NewSpacer(),
		widget.NewLabel("Port:"),
		portEntry,
		widget.NewLabel("Access Token:"),
		container.NewBorder(nil, nil, nil, container.NewHBox(regenBtn, copyTokenBtn), tokenEntry),
		layout.NewSpacer(),
		saveBtn,
	)

	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()

	return finalConfig, isSaved
}

func generateNewToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// FIXED: Outbound connection method to find the real IP
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1" // Fallback if no internet
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
