// Package copyfile provides an implementation of copying files, analogous to
// the cp command and some of its flags.
//
// The most basic mode is a physical byte-by-byte copy of the file. Optionally
// it can attempt to hardlink files instead of copying (falling back to a copy on
// failure); this is a bit like cp --link.
// It can also take advantage of copy-on-write semantics on some platforms,
// similar to cp --reflink. Currently this is supported only on Linux and OSX,
// and requires an underlying filesystem supporting it (e.g. btrfs, bcachefs, zfs, etc).
//
// When performing a physical file copy, they are created with a temporary name and moved
// into place once written, to try to avoid others seeing a partially-written file.
package copyfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"syscall"
)

// A Copier performs copying of files with some policies set.
// The zero Copier is safe to use and comes with default options.
//
// By default it attempts to clone files via copy-on-write but will
// detect if this appears to be unsupported and stop attempting it.
type Copier struct {
	// True to disable copy-on-write attempts via FICLONE on Linux.
	// N.B. This does not affect OSX where it is not explicitly selectable.
	DisableCOW bool
	// True to always force copy-on-write and never attempt to disable it.
	AlwaysCOW bool
	// True to always select a generic implementation (never use anything
	// platform-specific). This will be forced if compiling without cgo.
	Generic bool
	// True if we've tried copy-on-write and it was unsupported.
	cowFailed bool
}

// Copy copies a file from src to dest, following the policies on the Copier.
func (c *Copier) Copy(src, dest string) error {
	return c.CopyMode(src, dest, 0644)
}

// CopyMode copies a file from src to dest, setting the given file mode on dest.
func (c *Copier) CopyMode(src, dest string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	if !c.cowFailed && !c.Generic {
		return c.copySpecialised(srcFile, dest, mode)
	}
	return c.WriteFile(srcFile, dest, mode)
}

// WriteFile writes data from a reader to the file 'dest', with an attempt to perform
// a copy & rename to avoid chaos if anything goes wrong partway.
func (c *Copier) WriteFile(r io.Reader, dest string, mode os.FileMode) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	dir, file := path.Split(dest)
	tempFile, err := ioutil.TempFile(dir, file)
	if err != nil {
		return err
	} else if _, err := io.Copy(tempFile, r); err != nil {
		os.Remove(tempFile.Name()) // From here this is best-effort only
		return err
	} else if err := tempFile.Close(); err != nil {
		os.Remove(tempFile.Name())
		return err
	} else if err := os.Chmod(tempFile.Name(), mode); err != nil {
		os.Remove(tempFile.Name())
		return err
	}
	return os.Rename(tempFile.Name(), dest)
}

// Link hard-links a file from src to dest, falling back to a copy if the link fails.
func (c *Copier) Link(src, dest string) error {
	return c.LinkMode(src, dest, 0644)
}

// LinkMode hard-links a file from src to dest, falling back to a copy with the given mode if the link fails.
// The mode obviously does not apply if the file is linked.
func (c *Copier) LinkMode(src, dest string, mode os.FileMode) error {
	if err := os.Link(src, dest); err == nil {
		return nil
	} else if runtime.GOOS != "linux" && os.IsNotExist(err) {
		// There is an awkward issue on several non-Linux platforms where links to
		// symlinks actually try to link to the target rather than the link itself.
		// In that case we try to recreate a similar symlink at the destination.
		if info, err := os.Lstat(src); err == nil && (info.Mode()&os.ModeSymlink) != 0 {
			link, err := os.Readlink(src)
			if err != nil {
				return err
			}
			return os.Symlink(link, dest)
		}
		return err
	}
	return c.CopyMode(src, dest, mode)
}

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
