package mutex

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_ComplexMutexRun(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	ComplexMutexRun()

	w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "34320") {
		t.Error("Unexpected result")
	}
}
