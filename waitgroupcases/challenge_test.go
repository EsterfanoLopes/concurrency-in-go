package waitgroupcases

import (
	"io"
	"os"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	modifiedMsgValue := "modified value"

	var wg sync.WaitGroup
	wg.Add(unity)
	updateMessage(modifiedMsgValue, &wg)

	if msg != modifiedMsgValue {
		t.Errorf("unexpected value for msg. expeceted %s, found %s", modifiedMsgValue, msg)
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printMessage()

	w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if output != "\n" {
		t.Errorf("unexpected result. received %s. expected %s", output, msg)
	}
}

func Test_Challenge(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	Challenge()

	w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if output != "Hello, universe!\nHello, cosmos!\nHello, world!\n" {
		t.Errorf("unexpected result. received %s. expected %s", output, msg)
	}
}
