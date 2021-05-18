package dcmd

import (
	"bufio"
	"fmt"
	"github.com/osgochina/donkeygo/text/dstr"
	"os"
)

// Scan 交互操作
func Scan(info ...interface{}) string {
	fmt.Print(info...)
	return readline()
}

// Scanf  支持格式化参数的交互操作
func Scanf(format string, info ...interface{}) string {
	fmt.Printf(format, info...)
	return readline()
}

//从标准输入读入一行
func readline() string {
	var s string
	buf := bufio.NewReader(os.Stdin)
	s, _ = buf.ReadString('\n')
	s = dstr.Trim(s)
	return s
}
