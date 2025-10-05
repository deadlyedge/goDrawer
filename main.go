package main

import (
	"log"

	"github.com/deadlyedge/goDrawer/internal/settings"
	"github.com/deadlyedge/goDrawer/internal/ui"
)

func main() {
	// Read and print settings first
	appSettings, err := settings.Read("drawers-settings.toml")
	if err != nil {
		log.Printf("failed to read settings: %v", err)
	} else {
		settings.Print(appSettings)
	}

	// Run the UI window
	ui.RunWindow()
}
