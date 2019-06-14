//+build !cgo

package copyfile

import "os"

// This is a generic implementation that just calls through to WriteFile.
func (c *Copier) copySpecialised(srcFile *os.File, dest string, mode os.FileMode) error {
	return c.WriteFile(srcFile, dest, mode)
}
