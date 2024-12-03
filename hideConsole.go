package main

import (
	"syscall"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	user32               = syscall.NewLazyDLL("user32.dll")
	procShowWindow       = user32.NewProc("ShowWindow")
)

const (
	SW_HIDE = 0
	SW_SHOW = 5
)

func hideConsole() {
	hwnd, _, err := procGetConsoleWindow.Call()
	if hwnd != 0 && err == nil {
		procShowWindow.Call(hwnd, SW_HIDE)
	}
}

func showConsole() {
	hwnd, _, err := procGetConsoleWindow.Call()
	if hwnd != 0 && err == nil {
		procShowWindow.Call(hwnd, SW_SHOW)
	}
}
