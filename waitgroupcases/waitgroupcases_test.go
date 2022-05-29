package waitgroupcases

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printString(t *testing.T) {
	testString := "something"
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)

	go printString(testString, &wg)

	wg.Wait()

	w.Close()

	result, _ := io.ReadAll(r)

	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, testString) {
		t.Errorf("Unexpected result for output. %s", output)
	}
}
