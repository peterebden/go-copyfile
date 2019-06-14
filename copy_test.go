package copyfile

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// temp until I get an internet connection and can set up testify
func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func assertEqual(t *testing.T, b1, b2 []byte) {
	if !bytes.Equal(b1, b2) {
		t.Errorf("Unexpectedly not equal: %s %s", b1, b2)
	}
}

func assertNotEqual(t *testing.T, b1, b2 []byte) {
	if bytes.Equal(b1, b2) {
		t.Errorf("Unexpectedly equal: %s %s", b1, b2)
	}
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Errorf("Unexpectedly false")
	}
}

func assertFalse(t *testing.T, b bool) {
	if b {
		t.Errorf("Unexpectedly true")
	}
}

func TestCopy(t *testing.T) {
	b, err := ioutil.ReadFile("test_data/test.txt")
	assertNoError(t, err)
	var c Copier
	err = c.Copy("test_data/test.txt", "test_data/test2.txt")
	assertNoError(t, err)
	b2, err := ioutil.ReadFile("test_data/test2.txt")
	assertNoError(t, err)
	assertEqual(t, b, b2)
	// Write some stuff to the new file, it should not modify the original under
	// any circumstances.
	err = ioutil.WriteFile("test_data/test2.txt", []byte("testing"), 0644)
	assertNoError(t, err)
	b, err = ioutil.ReadFile("test_data/test.txt")
	assertNoError(t, err)
	assertNotEqual(t, []byte("testing"), b)
}

func TestIsSameFile(t *testing.T) {
	var c Copier
	assertTrue(t, c.IsSameFile("test_data/test.txt", "test_data/test.txt"))
	assertFalse(t, c.IsSameFile("test_data/test.txt", "test_data"))
}
