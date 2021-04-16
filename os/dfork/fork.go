package dfork

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"syscall"
)

var forkFuncMap = make(map[string]func())

// AddMethod 添加要执行的方法
func AddMethod(name string, runFunc func()) {
	_, found := forkFuncMap[name]
	if found {
		panic(fmt.Sprintf("AddMethod func already registered under name %q", name))
	}
	forkFuncMap[name] = runFunc
}

// Run 运行
func Run() bool {
	runFunc, exists := forkFuncMap[os.Args[0]]
	if exists {
		runFunc()
		return true
	}
	return false
}

// Command 执行命令
func Command(args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: "/proc/self/exe",
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: unix.SIGTERM,
		},
	}
}
