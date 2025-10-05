package main

import (
	"fmt"
	"log"
	"math"
	"syscall"

	"github.com/lxn/walk"
	WD "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	targetAlpha byte   = 200 // 0-255 where 255 is fully opaque
	lwaAlpha    uint32 = 0x2
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
)

func main() {
	var (
		mainWindow      *walk.MainWindow
		alphaSlider     *walk.Slider
		alphaValueLabel *walk.Label
	)

	dragHandler := func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton && mainWindow != nil {
			beginWindowDrag(mainWindow)
		}
	}

	if err := (WD.MainWindow{
		AssignTo:    &mainWindow,
		Title:       "Semi-transparent Test Window",
		MinSize:     WD.Size{Width: 400, Height: 300},
		Size:        WD.Size{Width: 400, Height: 300},
		Layout:      WD.VBox{MarginsZero: true, Spacing: 12},
		OnMouseDown: dragHandler,
		Children: []WD.Widget{
			WD.Composite{
				Layout:      WD.VBox{Margins: WD.Margins{Left: 12, Top: 12, Right: 12, Bottom: 12}, Spacing: 8},
				OnMouseDown: dragHandler,
				Children: []WD.Widget{
					WD.Label{
						Text:        "This window uses Walk and Win32 layered styles for alpha blending.",
						OnMouseDown: dragHandler,
					},
					WD.Label{
						AssignTo:    &alphaValueLabel,
						Text:        opacityLabelText(int(targetAlpha)),
						OnMouseDown: dragHandler,
					},
					WD.Slider{
						AssignTo: &alphaSlider,
						MinValue: 50,
						MaxValue: 255,
						Tracking: true,
						OnValueChanged: func() {
							if alphaSlider == nil {
								return
							}

							current := alphaSlider.Value()

							if alphaValueLabel != nil {
								alphaValueLabel.SetText(opacityLabelText(current))
							}

							if mainWindow != nil {
								if err := makeWindowSemiTransparent(mainWindow, byte(current)); err != nil {
									log.Printf("failed to update window transparency: %v", err)
								}
							}
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatalf("failed to create main window: %v", err)
	}

	if err := makeWindowBorderless(mainWindow); err != nil {
		log.Fatalf("failed to remove window border: %v", err)
	}

	if alphaSlider != nil {
		alphaSlider.SetValue(int(targetAlpha))
	}

	if err := makeWindowSemiTransparent(mainWindow, targetAlpha); err != nil {
		log.Fatalf("failed to update window transparency: %v", err)
	}

	mainWindow.Run()
}

func makeWindowBorderless(mw *walk.MainWindow) error {
	hwnd := mw.Handle()
	if hwnd == 0 {
		return syscall.EINVAL
	}

	style := uint32(win.GetWindowLong(hwnd, win.GWL_STYLE))
	style &^= uint32(win.WS_CAPTION | win.WS_THICKFRAME | win.WS_MINIMIZE | win.WS_MAXIMIZE | win.WS_SYSMENU)
	style |= uint32(win.WS_POPUP)

	win.SetLastError(0)
	if prev := win.SetWindowLong(hwnd, win.GWL_STYLE, int32(style)); prev == 0 {
		if err := syscall.GetLastError(); err != nil && err != syscall.Errno(0) {
			return err
		}
	}

	if !win.SetWindowPos(hwnd, 0, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOZORDER|win.SWP_FRAMECHANGED|win.SWP_NOACTIVATE) {
		return syscall.GetLastError()
	}

	return nil
}

func makeWindowSemiTransparent(mw *walk.MainWindow, alpha byte) error {
	hwnd := mw.Handle()
	if hwnd == 0 {
		return syscall.EINVAL
	}

	exStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, exStyle|win.WS_EX_LAYERED)

	return setLayeredWindowAttributes(hwnd, win.COLORREF(0), alpha, lwaAlpha)
}

func setLayeredWindowAttributes(hwnd win.HWND, key win.COLORREF, alpha byte, flags uint32) error {
	r1, _, err := procSetLayeredWindowAttributes.Call(
		uintptr(hwnd),
		uintptr(key),
		uintptr(alpha),
		uintptr(flags),
	)
	if r1 == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return errno
		}
		return syscall.EINVAL
	}

	return nil
}

func beginWindowDrag(mw *walk.MainWindow) {
	hwnd := mw.Handle()
	if hwnd == 0 {
		return
	}

	win.ReleaseCapture()
	win.SendMessage(hwnd, win.WM_NCLBUTTONDOWN, uintptr(win.HTCAPTION), 0)
}

func opacityLabelText(alpha int) string {
	percent := int(math.Round(float64(alpha) * 100 / 255))
	return fmt.Sprintf("Opacity: %d (%d%%)", alpha, percent)
}
