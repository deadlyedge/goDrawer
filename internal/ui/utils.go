package ui

import (
	"fmt"
	"math"
	"syscall"
	// "unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	gdi32                          = syscall.NewLazyDLL("gdi32.dll")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procCreateRoundRectRgn         = gdi32.NewProc("CreateRoundRectRgn")
	procSetWindowRgn               = user32.NewProc("SetWindowRgn")
	procDeleteObject               = gdi32.NewProc("DeleteObject")
)

func makeWindowBorderless(mainWindow *walk.MainWindow) error {
	hwnd := mainWindow.Handle()
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

func makeWindowRoundedCorner(mw *walk.MainWindow, radius int) error {
	hwnd := mw.Handle()
	if hwnd == 0 {
		return syscall.EINVAL
	}

	// 获取窗口尺寸
	var rect win.RECT
	if !win.GetWindowRect(hwnd, &rect) {
		return syscall.GetLastError()
	}

	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)

	if width <= 0 || height <= 0 {
		return syscall.EINVAL
	}

	// 创建圆角区域
	r1, _, err := procCreateRoundRectRgn.Call(
		uintptr(0),      // nLeftRect
		uintptr(0),      // nTopRect
		uintptr(width),  // nRightRect
		uintptr(height), // nBottomRect
		uintptr(radius), // nWidthEllipse
		uintptr(radius), // nHeightEllipse
	)
	if r1 == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return errno
		}
		return syscall.EINVAL
	}
	rgn := win.HRGN(r1)

	// 应用区域到窗口
	r2, _, err := procSetWindowRgn.Call(
		uintptr(hwnd), // hWnd
		uintptr(rgn),  // hRgn
		uintptr(1),    // bRedraw
	)
	if r2 == 0 {
		// 清理区域对象
		procDeleteObject.Call(uintptr(rgn))
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return errno
		}
		return syscall.EINVAL
	}

	// 清理区域对象（成功应用后也需要清理）
	procDeleteObject.Call(uintptr(rgn))

	return nil
}

func makeWindowSemiTransparent(mw *walk.MainWindow, alpha byte) error {
	hwnd := mw.Handle()
	if hwnd == 0 {
		return syscall.EINVAL
	}

	exStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, exStyle|win.WS_EX_LAYERED)

	return setLayeredWindowAttributes(hwnd, win.COLORREF(0), alpha, uint32(lwaAlpha))
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
