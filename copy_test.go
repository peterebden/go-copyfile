package copyfile

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	b, err := ioutil.ReadFile("test_data/test.txt")
	assert.NoError(t, err)
}
