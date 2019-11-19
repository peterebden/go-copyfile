// +build !windows

package copyfile

import (
	"fmt"
	"os"
	"syscall"
)

// IsSameFile returns true if the two given paths refer to the same file.
func (c *Copier) IsSameFile(file1, file2 string) bool {
	i1, err1 := c.getInode(file1)
	i2, err2 := c.getInode(file2)
	return err1 == nil && err2 == nil && i1 == i2
}

func (c *Copier) getInode(filename string) (uint64, error) {
	fi, err := os.Lstat(filename)
	if err != nil {
		return 0, err
	}
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("Not a syscall.Stat_t")
	}
	return uint64(s.Ino), nil
}
