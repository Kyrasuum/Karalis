package rlx

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// InitWindow - Initialize window and OpenGL context
func InitWindow(width int32, height int32, title string) {
	Do(func() {
		rl.InitWindow(width, height, title)
	})
}

// CloseWindow - Close window and unload OpenGL context
func CloseWindow() {
	Do(func() {
		rl.CloseWindow()
	})
}

func SetConfigFlags(flags uint32) {
	Do(func() {
		rl.SetConfigFlags(flags)
	})
}

// WindowShouldClose - Check if application should close
func WindowShouldClose() bool {
	return Call(func() bool {
		return rl.WindowShouldClose()
	})
}

func IsWindowReady() bool {
	return Call(func() bool {
		return rl.IsWindowReady()
	})
}

func IsWindowFullscreen() bool {
	return Call(func() bool {
		return rl.IsWindowFullscreen()
	})
}

func IsWindowHidden() bool {
	return Call(func() bool {
		return rl.IsWindowHidden()
	})
}

func IsWindowMinimized() bool {
	return Call(func() bool {
		return rl.IsWindowMinimized()
	})
}

func IsWindowMaximized() bool {
	return Call(func() bool {
		return rl.IsWindowMaximized()
	})
}

func IsWindowFocused() bool {
	return Call(func() bool {
		return rl.IsWindowFocused()
	})
}

func IsWindowResized() bool {
	return Call(func() bool {
		return rl.IsWindowResized()
	})
}

func IsWindowState(flag uint32) bool {
	return Call(func() bool {
		return rl.IsWindowState(flag)
	})
}

func SetWindowState(flags uint32) {
	Do(func() {
		rl.SetWindowState(flags)
	})
}

func ClearWindowState(flags uint32) {
	Do(func() {
		rl.ClearWindowState(flags)
	})
}

func ToggleFullscreen() {
	Do(func() {
		rl.ToggleFullscreen()
	})
}

func ToggleBorderlessWindowed() {
	Do(func() {
		rl.ToggleBorderlessWindowed()
	})
}

func MaximizeWindow() {
	Do(func() {
		rl.MaximizeWindow()
	})
}

func MinimizeWindow() {
	Do(func() {
		rl.MinimizeWindow()
	})
}

func RestoreWindow() {
	Do(func() {
		rl.RestoreWindow()
	})
}

func SetWindowIcon(image rl.Image) {
	Do(func() {
		rl.SetWindowIcon(image)
	})
}

func SetWindowIcons(images []rl.Image, count int32) {
	copied := append([]rl.Image(nil), images...)
	Do(func() {
		rl.SetWindowIcons(copied, count)
	})
}

func SetWindowTitle(title string) {
	Do(func() {
		rl.SetWindowTitle(title)
	})
}

func SetWindowPosition(x int, y int) {
	Do(func() {
		rl.SetWindowPosition(x, y)
	})
}

func SetWindowMonitor(monitor int) {
	Do(func() {
		rl.SetWindowMonitor(monitor)
	})
}

func SetWindowMinSize(width int, height int) {
	Do(func() {
		rl.SetWindowMinSize(width, height)
	})
}

func SetWindowMaxSize(width int, height int) {
	Do(func() {
		rl.SetWindowMaxSize(width, height)
	})
}

func SetWindowSize(width int, height int) {
	Do(func() {
		rl.SetWindowSize(width, height)
	})
}

func SetWindowOpacity(opacity float32) {
	Do(func() {
		rl.SetWindowOpacity(opacity)
	})
}

func SetWindowFocused() {
	Do(func() {
		rl.SetWindowFocused()
	})
}

func GetWindowHandle() unsafe.Pointer {
	return Call(func() unsafe.Pointer {
		return rl.GetWindowHandle()
	})
}

func GetScreenWidth() int {
	return Call(func() int {
		return rl.GetScreenWidth()
	})
}

func GetScreenHeight() int {
	return Call(func() int {
		return rl.GetScreenHeight()
	})
}

func GetRenderWidth() int {
	return Call(func() int {
		return rl.GetRenderWidth()
	})
}

func GetRenderHeight() int {
	return Call(func() int {
		return rl.GetRenderHeight()
	})
}

func GetMonitorCount() int {
	return Call(func() int {
		return rl.GetMonitorCount()
	})
}

func GetCurrentMonitor() int {
	return Call(func() int {
		return rl.GetCurrentMonitor()
	})
}

func GetMonitorPosition(monitor int) rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetMonitorPosition(monitor)
	})
}

func GetMonitorWidth(monitor int) int {
	return Call(func() int {
		return rl.GetMonitorWidth(monitor)
	})
}

func GetMonitorHeight(monitor int) int {
	return Call(func() int {
		return rl.GetMonitorHeight(monitor)
	})
}

func GetMonitorPhysicalWidth(monitor int) int {
	return Call(func() int {
		return rl.GetMonitorPhysicalWidth(monitor)
	})
}

func GetMonitorPhysicalHeight(monitor int) int {
	return Call(func() int {
		return rl.GetMonitorPhysicalHeight(monitor)
	})
}

func GetMonitorRefreshRate(monitor int) int {
	return Call(func() int {
		return rl.GetMonitorRefreshRate(monitor)
	})
}

func GetWindowPosition() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetWindowPosition()
	})
}

func GetWindowScaleDPI() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetWindowScaleDPI()
	})
}

func GetMonitorName(monitor int) string {
	return Call(func() string {
		return rl.GetMonitorName(monitor)
	})
}

func SetClipboardText(text string) {
	Do(func() {
		rl.SetClipboardText(text)
	})
}

func GetClipboardText() string {
	return Call(func() string {
		return rl.GetClipboardText()
	})
}

func GetClipboardImage() rl.Image {
	return Call(func() rl.Image {
		return rl.GetClipboardImage()
	})
}

func EnableEventWaiting() {
	Do(func() {
		rl.EnableEventWaiting()
	})
}

func DisableEventWaiting() {
	Do(func() {
		rl.DisableEventWaiting()
	})
}

func ShowCursor() {
	Do(func() {
		rl.ShowCursor()
	})
}

func HideCursor() {
	Do(func() {
		rl.HideCursor()
	})
}

func IsCursorHidden() bool {
	return Call(func() bool {
		return rl.IsCursorHidden()
	})
}

func EnableCursor() {
	Do(func() {
		rl.EnableCursor()
	})
}

func DisableCursor() {
	Do(func() {
		rl.DisableCursor()
	})
}

func IsCursorOnScreen() bool {
	return Call(func() bool {
		return rl.IsCursorOnScreen()
	})
}
