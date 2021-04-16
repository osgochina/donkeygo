package dfork_test

import (
	"donkeygo/os/dfork"
	"testing"
)
import "github.com/gogf/gf/test/gtest"

func init() {
	dfork.AddMethod("fork", func() {
		panic("Return Error")
	})
	dfork.Run()
}

func TestAddMethod(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			gtest.AssertEQ(`AddMethod func already registered under name "fork"`, r)
		}
	}()
	dfork.AddMethod("fork", func() {})
}

func TestCommand(t *testing.T) {
	cmd := dfork.Command("fork")
	w, err := cmd.StdinPipe()
	gtest.AssertNil(err)
	defer w.Close()

	err = cmd.Start()
	gtest.AssertNil(err)
	err = cmd.Wait()
	gtest.AssertNil(err)
}
