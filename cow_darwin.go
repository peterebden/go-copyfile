//+build cgo

package copyfile

import (
	"os"
	"syscall"
)

// #include <errno.h>
// #include <sys/ioctl.h>
import "C"

func (c *Copier) copySpecialised(srcFile *os.File, dest string, mode os.FileMode) error {
	return C.copyfile(srcFile, dest)
}
