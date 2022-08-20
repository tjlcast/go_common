package file_utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestReadFileLines(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := strings.Join([]string{dir, "file_helper_test.go"}, string(os.PathSeparator))
	lines, err := ReadFileLines(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(lines))

	// first line.
	if lines[0] != "package file_utils" {
		assert.Fail(t, "First line is not right.")
	}
	// last line.
	if lines[len(lines)-1] != "// end for assert" {
		assert.Fail(t, "Last line is not right.")
	}
}
// end for assert