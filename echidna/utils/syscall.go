package utils

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

var (
	GetThreadPID = user32.NewProc("GetWindowThreadProcessId")
	SendMessage  = user32.NewProc("SendMessageW")
	OpenProcess  = kernel32.NewProc("OpenProcess")
	CloseHandle  = kernel32.NewProc("CloseHandle")
	VirtualAlloc = kernel32.NewProc("VirtualAllocEx")
	VirtualFree  = kernel32.NewProc("VirtualFreeEx")
	WriteMemory  = kernel32.NewProc("WriteProcessMemory")
	ReadMemory   = kernel32.NewProc("ReadProcessMemory")

	findWindowEx = user32.NewProc("FindWindowExW")
	findWindow   = user32.NewProc("FindWindowW")
	messageBox   = user32.NewProc("MessageBoxW")
)

func FindWindowEx(parent uintptr, class string) uintptr {
	var classPtr, _ = syscall.UTF16PtrFromString(class)

	handle, _, _ := findWindowEx.Call(
		parent,
		0,
		uintptr(unsafe.Pointer(classPtr)),
		0,
	)

	return handle
}

func FindWindow(class string, title string) uintptr {
	var classPtr, _ = syscall.UTF16PtrFromString(class)
	var titlePtr, _ = syscall.UTF16PtrFromString(title)

	handle, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(classPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
	)

	return handle
}

func Messagebox(message string) {
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	messageBox.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(0),
		uintptr(0x00000000|0x00040000|0x00001000),
	)
}
