package copyfile

import (
	"os"
)

// IsSameFile returns true if the two given paths refer to the same file.
func (c *Copier) IsSameFile(file1, file2 string) bool {
	fi1, err := os.Stat(file1)
	if err != nil {
		return false
	}

	fi2, err := os.Stat(file2)
	if err != nil {
		return false
	}

	return os.SameFile(fi1, fi2)
}
