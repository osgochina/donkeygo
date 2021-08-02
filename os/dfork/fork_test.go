package dfork_test

import (
	"github.com/osgochina/donkeygo/os/dfork"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func init() {
	dfork.AddMethod("fork", func() {
		panic("Return Error")
	})
	dfork.Run()
}

func TestAddMethod(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			dtest.AssertEQ(`AddMethod func already registered under name "fork"`, r)
		}
	}()
	dfork.AddMethod("fork", func() {})
}

func TestCommand(t *testing.T) {
	cmd := dfork.Command("fork")
	w, err := cmd.StdinPipe()
	dtest.AssertNil(err)
	defer w.Close()

	err = cmd.Start()
	dtest.AssertNil(err)
	err = cmd.Wait()
	dtest.AssertNil(err)
}
