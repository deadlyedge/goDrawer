package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/deadlyedge/goDrawer/internal/settings"
)

// App struct
type App struct {
	ctx context.Context
}

var path = "drawers-settings.toml"
var appSettings, _ = settings.Read(path)

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.WindowSetPosition(
		ctx,
		appSettings.WindowPosition.X,
		appSettings.WindowPosition.Y,
	)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Save window position
func (a *App) SaveWindowPosition(x int, y int) {
	appSettings.WindowPosition.X = x
	appSettings.WindowPosition.Y = y
	settings.Update(path, appSettings)
}
// v2: BeforeClose 钩子
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
    x, y := runtime.WindowGetPosition(ctx)
		a.SaveWindowPosition(x, y)
    return false // 允许关闭
}