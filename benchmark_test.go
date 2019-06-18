package copyfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var dir string
var file string

func BenchmarkCopy(b *testing.B) {
	var c Copier
	for i := 0; i < b.N; i++ {
		out := fmt.Sprintf(path.Join(dir, "out_%d"), i)
		if err := c.Copy(file, out); err != nil {
			panic(err)
		}
		os.Remove(out)
	}
}

func BenchmarkLink(b *testing.B) {
	var c Copier
	for i := 0; i < b.N; i++ {
		out := fmt.Sprintf(path.Join(dir, "out_%d"), i)
		if err := c.Link(file, out); err != nil {
			panic(err)
		}
		os.Remove(out)
	}
}

func BenchmarkCopyGeneric(b *testing.B) {
	c := Copier{Generic: true}
	for i := 0; i < b.N; i++ {
		out := fmt.Sprintf(path.Join(dir, "out_%d"), i)
		if err := c.Copy(file, out); err != nil {
			panic(err)
		}
		os.Remove(out)
	}
}

func testMain(m *testing.M) int {
	dir, _ = ioutil.TempDir("", "go-copyfile")
	defer os.RemoveAll(dir)
	f, err := ioutil.TempFile(dir, "ref")
	if err != nil {
		panic(err)
	}
	file = f.Name()
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}
