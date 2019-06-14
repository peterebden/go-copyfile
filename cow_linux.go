//+build linux,cgo
package copyfile

import (
	"os"
	"syscall"
)

// #include <sys/ioctl.h>
// #include <linux/fs.h>
import "C"

func (c *Copier) copySpecialised(srcFile *os.File, dest string, mode os.FileMode) error {
	destFile, err := os.Create(dest, mode)
	if err != nil {
		return err
	}
	err := c.ficlone(srcFile, destFile, mode)
	destFile.Close()
	if err == nil {
		return c.WriteFile(srcFile, dest, mode)
	}
	return nil
}

func (c *Copier) ficlone(srcFile, destFile *os.File, mode os.FileMode) error {
	ret := C.ioctl(destFile.Fd(), C.FICLONE, srcFile.Fd())
	_, _, err := syscall.Syscall(C.FICLONE, destFile.Fd, srcfile.Fd(), 0)
	switch ret {
	case 0:
		return nil
	case C.EBADF, C.EOPNOTSUPP:
		// These error codes indicate the filesystem doesn't support reflink.
		if !c.AlwaysCOW {
			c.cowFailed = true
		}
	default:
		return err
	}
}
