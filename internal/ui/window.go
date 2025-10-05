package ui

import (
	"log"

	"github.com/deadlyedge/goDrawer/internal/settings"
	"github.com/lxn/walk"
	WD "github.com/lxn/walk/declarative"
)

const (
	targetAlpha byte   = 200 // 0-255 where 255 is fully opaque
	lwaAlpha    uint32 = 0x2
)

func RunWindow() {
	var (
		mainWindow      *walk.MainWindow
		settingsButton  *walk.PushButton
		addDrawerButton *walk.PushButton
		drawersList     *walk.ListBox
	)

	// 读取设置
	appSettings, err := settings.Read("drawers-settings.toml")
	if err != nil {
		log.Printf("failed to read settings: %v", err)
		appSettings = &settings.Settings{Drawers: []settings.Drawer{}}
	}

	dragHandler := func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton && mainWindow != nil {
			beginWindowDrag(mainWindow)
		}
	}

	if err := (WD.MainWindow{
		AssignTo:    &mainWindow,
		Title:       "goDrawer",
		MinSize:     WD.Size{Width: 300, Height: 400},
		Size:        WD.Size{Width: 350, Height: 450},
		Layout:      WD.VBox{MarginsZero: true, Spacing: 8},
		// OnMouseDown: dragHandler,
		Children: []WD.Widget{
			// 第一行：拖动区域 + 设置按钮
			WD.Composite{
				Layout: WD.HBox{Spacing: 0},
				Children: []WD.Widget{
					// 拖动区域（左边占位符）
					WD.Composite{
						Layout:   WD.HBox{},
						AssignTo: nil,
					},
					// 设置按钮右侧对齐
					WD.Composite{
						Layout: WD.HBox{},
						Children: []WD.Widget{
							WD.PushButton{
								AssignTo: &settingsButton,
								Text:     "⚙️",
								MinSize:  WD.Size{Width: 30, Height: 30},
								MaxSize:  WD.Size{Width: 30, Height: 30},
								OnClicked: func() {
									// TODO: 打开设置面板
									log.Println("Settings button clicked")
								},
							},
						},
					},
				},
				OnMouseDown: dragHandler,
			},
			// 第二行：抽屉列表
			WD.ListBox{
				AssignTo:    &drawersList,
				Model:       getDrawerNames(appSettings.Drawers),
				OnMouseDown: dragHandler,
			},
			// 第三行：添加抽屉按钮
			WD.PushButton{
				AssignTo: &addDrawerButton,
				Text:     "+ 添加抽屉",
				OnClicked: func() {
					// TODO: 添加抽屉功能
					log.Println("Add drawer button clicked")
				},
			},
		},
	}).Create(); err != nil {
		log.Fatalf("failed to create main window: %v", err)
	}

	if err := makeWindowBorderless(mainWindow); err != nil {
		log.Fatalf("failed to remove window border: %v", err)
	}

	// 设置半透明
	if err := makeWindowSemiTransparent(mainWindow, 240); err != nil {
		log.Fatalf("failed to update window transparency: %v", err)
	}

	// 设置圆角
	if err := makeWindowRoundedCorner(mainWindow, 15); err != nil {
		log.Printf("failed to set window corners: %v", err)
	}

	mainWindow.Run()
}

func getDrawerNames(drawers []settings.Drawer) []string {
	names := make([]string, len(drawers))
	for i, drawer := range drawers {
		names[i] = drawer.Name
	}
	return names
}
