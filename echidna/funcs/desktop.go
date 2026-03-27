package funcs

import (
	"syscall"
	"unsafe"

	"echidna/store"
	"echidna/utils"
)

const (
	pPROCESS_VM_OPERATION = 0x0008
	pPROCESS_VM_READ      = 0x0010
	pPROCESS_VM_WRITE     = 0x0020

	pLVM_GETITEMCOUNT    = 0x1004
	pLVM_GETITEMPOSITION = 0x1010
	pLVM_GETITEMW        = 0x104B
	pLVIF_TEXT           = 0x0001

	pMEM_COMMIT     = 0x1000
	pMEM_RELEASE    = 0x8000
	pPAGE_READWRITE = 0x04

	pBUFFER_SIZE = 256
)

type listViewItem struct {
	Mask       uint32
	IItem      int32
	ISubItem   int32
	State      uint32
	StateMask  uint32
	PszText    uintptr
	CchTextMax int32
	IImage     int32
	LParam     uintptr
	IIndent    int32
	IGroupId   int32
	CColumns   uint32
	PuColumns  uintptr
	PiColFmt   uintptr
	IGroup     int32
}

type point struct {
	X int32
	Y int32
}

func GetDesktopIcons() []store.DesktopIcon {
	progman := utils.FindWindow("Progman", "Program Manager")
	shell := utils.FindWindowEx(progman, "SHELLDLL_DefView")

	view := utils.FindWindowEx(shell, "SysListView32")
	if view == 0 {
		return nil
	}

	var pid uint32
	utils.GetThreadPID.Call(view, uintptr(unsafe.Pointer(&pid)))

	proc, _, _ := utils.OpenProcess.Call(
		pPROCESS_VM_OPERATION|pPROCESS_VM_READ|pPROCESS_VM_WRITE,
		0,
		uintptr(pid),
	)

	if proc == 0 {
		return nil
	}

	defer utils.CloseHandle.Call(proc)

	count, _, _ := utils.SendMessage.Call(view, pLVM_GETITEMCOUNT, 0, 0)

	items := int(count)
	if items == 0 {
		return nil
	}

	vName, _, _ := utils.VirtualAlloc.Call(
		proc,
		0,
		pBUFFER_SIZE*2,
		pMEM_COMMIT,
		pPAGE_READWRITE,
	)

	vItem, _, _ := utils.VirtualAlloc.Call(
		proc,
		0,
		unsafe.Sizeof(listViewItem{}),
		pMEM_COMMIT,
		pPAGE_READWRITE,
	)

	vPos, _, _ := utils.VirtualAlloc.Call(
		proc,
		0,
		unsafe.Sizeof(point{}),
		pMEM_COMMIT,
		pPAGE_READWRITE,
	)

	defer utils.VirtualFree.Call(proc, vName, 0, pMEM_RELEASE)
	defer utils.VirtualFree.Call(proc, vItem, 0, pMEM_RELEASE)
	defer utils.VirtualFree.Call(proc, vPos, 0, pMEM_RELEASE)

	icons := make([]store.DesktopIcon, 0, items)

	for i := range items {
		item := listViewItem{
			Mask:       pLVIF_TEXT,
			IItem:      int32(i),
			PszText:    vName,
			CchTextMax: pBUFFER_SIZE,
		}

		utils.WriteMemory.Call(
			proc,
			vItem,
			uintptr(unsafe.Pointer(&item)),
			unsafe.Sizeof(item), 0,
		)

		utils.SendMessage.Call(
			view,
			pLVM_GETITEMW,
			uintptr(i),
			vName,
		)

		var name [pBUFFER_SIZE]uint16
		utils.ReadMemory.Call(
			proc,
			vName,
			uintptr(unsafe.Pointer(&name[0])),
			pBUFFER_SIZE*2,
			0,
		)

		utils.SendMessage.Call(
			view,
			pLVM_GETITEMPOSITION,
			uintptr(i),
			vPos,
		)

		var position point
		utils.ReadMemory.Call(
			proc,
			vPos,
			uintptr(unsafe.Pointer(&position)),
			unsafe.Sizeof(position),
			0,
		)

		icons = append(icons, store.DesktopIcon{
			Name: syscall.UTF16ToString(name[:]),
			X:    int(position.X),
			Y:    int(position.Y),
		})
	}

	return icons
}
