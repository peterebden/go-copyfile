package copyfile

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32   = syscall.NewLazyDLL("kernel32.dll")
	procCopyFileW = modkernel32.NewProc("CopyFileW")
)

func (c *Copier) copySpecialised(srcFile *os.File, dest string, mode os.FileMode) error {
	var bFailIfExists uint32 = 1

	lpExistingFileName, err := syscall.UTF16PtrFromString(srcFile.Name())
	if err != nil {
		return err
	}
	lpNewFileName, err := syscall.UTF16PtrFromString(dest)
	if err != nil {
		return err
	}
	r1, _, err := syscall.Syscall(
		procCopyFileW.Addr(),
		3,
		uintptr(unsafe.Pointer(lpExistingFileName)),
		uintptr(unsafe.Pointer(lpNewFileName)),
		uintptr(bFailIfExists))
	if r1 == 0 {
		return fmt.Errorf("failed CopyFileW Win32 call from '%s' to '%s': %s", srcFile.Name(), dest, err)
	}
	return nil
}
