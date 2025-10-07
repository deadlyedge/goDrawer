package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	// "github.com/wailsapp/wails/v3/pkg/options"
	// "github.com/wailsapp/wails/v3/pkg/options/windows"
	// "github.com/wailsapp/wails/v3/pkg/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := application.New(application.Options{
		Name:        "goDrawer",
		Description: "A mini windows desktop drawer app",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})


	// Create new window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           "goDrawer",
		Width:           256,
		Height:          256,
		Frameless:       true,
		DisableResize:   true,
		BackgroundColour: application.NewRGBA(27, 38, 54, 200), // 半透明背景
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			BackdropType: application.Auto,
		},
		URL: "/",
	})

	// TODO: Add window position saving on close in Wails v3 compatible way

	// Create system tray
	systray := app.SystemTray.New()
	systray.SetTooltip("My Application Tooltip")

	// Create menu
	menu := app.NewMenu()
	menu.Add("hello world").OnClick(func(ctx *application.Context) {
		window.Show()
		log.Println("hello world")
	})
	// Add menu to system tray
	systray.SetMenu(menu)

	// Create application with options
	// err := wails.Run(&options.App{
	// 	Title: "goDrawer",
	// 	AssetServer: &assetserver.Options{
	// 		Assets: assets,
	// 	},
	// 	BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 100},
	// 	Bind: []interface{}{
	// 		app,
	// 	},
	// })
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
