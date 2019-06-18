//+build cgo

package copyfile

import (
	"os"
	"syscall"
)

// #include <copyfile.h>
import "C"

func (c *Copier) copySpecialised(srcFile *os.File, dest string, mode os.FileMode) error {
	destFile, err := os.OpenFile(dest, os.O_WRONLY | os.O_CREATE, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()
	state := C.copyfile_state_alloc()
	defer C.copyfile_state_free(state) 
	if err := C.fcopyfile(C.int(srcFile.Fd()), C.int(destFile.Fd()), state, C.COPYFILE_DATA); err != 0 {
		return syscall.Errno(err)
	}
	return nil
}
