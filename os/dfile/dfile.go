package dfile

import (
	"donkeygo/container/dtype"
	"donkeygo/text/dstr"
	"donkeygo/util/dconv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	// Separator 操作系统文件路径分隔符
	Separator = string(filepath.Separator)

	// DefaultPermOpen 是默认打开文件的权限值
	DefaultPermOpen = os.FileMode(0666)

	// DefaultPermCopy 是默认copy文件或文件夹的权限值
	DefaultPermCopy = os.FileMode(0777)

	//main包的路径
	mainPkgPath = dtype.NewString()

	//程序自身的路径
	selfPath = ""

	//操作系统临时目录路径
	tempDir = "/tmp"
)

func init() {

	//确认临时目录的路径
	if Separator != "/" || !Exists(tempDir) {
		tempDir = os.TempDir()
	}

	//初始化变脸 selfPath，确认当前运行程序的绝对地址
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

// Mkdir 根据<path>给定的路径递归创建目录，建议使用绝对路径
func Mkdir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Create 递归创建path给定的地址，建议使用绝对路径
func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

// Open 只读的方式打开文件或文件夹.
func Open(path string) (*os.File, error) {
	return os.Open(path)
}

// OpenFile opens file/directory with custom <flag> and <perm>.
// The parameter <flag> is like: O_RDONLY, O_RDWR, O_RDWR|O_CREATE|O_TRUNC, etc.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

// OpenWithFlagPerm opens file/directory with custom <flag> and <perm>.
// The parameter <flag> is like: O_RDONLY, O_RDWR, O_RDWR|O_CREATE|O_TRUNC, etc.
// The parameter <perm> is like: 0600, 0666, 0777, etc.
func OpenWithFlagPerm(path string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Join joins string array paths with file separator of current system.
func Join(paths ...string) string {
	var s string
	for _, path := range paths {
		if s != "" {
			s += Separator
		}
		s += dstr.TrimRight(path, Separator)
	}
	return s
}

// Exists 判断指定路径的文件是否存在
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir 判断path路径是否是目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Pwd 获取当前运行目录
func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

// Chdir 修改当前工作目录修改为 dir 指定的目录
func Chdir(dir string) error {
	return os.Chdir(dir)
}

// IsFile 判断传入的path是否是存在的文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Info Alias of Stat.
// See Stat.
func Info(path string) (os.FileInfo, error) {
	return Stat(path)
}

// Stat 获取文件的元数据
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Move 把scr指向的文件移动到dst
func Move(src string, dst string) error {
	return os.Rename(src, dst)
}

// Rename 修改文件名
func Rename(src string, dst string) error {
	return Move(src, dst)
}

// DirNames 返回path路径下的文件名列表。
// 注意返回的不是绝对路径
func DirNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func Glob(pattern string, onlyNames ...bool) ([]string, error) {
	if list, err := filepath.Glob(pattern); err == nil {
		if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
			array := make([]string, len(list))
			for k, v := range list {
				array[k] = Basename(v)
			}
			return array, nil
		}
		return list, nil
	} else {
		return nil, err
	}
}

// Remove deletes all file/directory with <path> parameter.
// If parameter <path> is directory, it deletes it recursively.
func Remove(path string) error {
	return os.RemoveAll(path)
}

// IsReadable 检查文件是否可读
func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, DefaultPermOpen)
	if err != nil {
		result = false
	}
	_ = file.Close()
	return result
}

// IsWritable 判断给定的路径是否可写，支持文件夹和文件
func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		// If it's a directory, create a temporary file to test whether it's writable.
		tmpFile := strings.TrimRight(path, Separator) + Separator + dconv.String(time.Now().UnixNano())
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			_ = f.Close()
			_ = Remove(tmpFile)
		}
	} else {
		// 如果是文件，那么判断文件是否可打开
		file, err := os.OpenFile(path, os.O_WRONLY, DefaultPermOpen)
		if err != nil {
			result = false
		}
		_ = file.Close()
	}
	return result
}

// Chmod See os.Chmod.
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// Abs 获取绝对路径
func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}

// RealPath 判断传入的path是否是真是存在的地址
func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

// SelfPath 获取当前可执行文件的地址
func SelfPath() string {
	return selfPath
}

// SelfName 获取当前运行程序的名字
func SelfName() string {
	return Basename(SelfPath())
}

// SelfDir 判断当前运行的目录
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename returns the last element of path, which contains file extension.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Basename returns a single separator.
// Example:
// /var/www/file.js -> file.js
// file.js          -> file.js
func Basename(path string) string {
	return filepath.Base(path)
}

// Name returns the last element of path without file extension.
// Example:
// /var/www/file.js -> file
// file.js          -> file
func Name(path string) string {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

//Dir 返回Path的最后一个元素以外的所有元素，通常是路径的目录。
//删除最后一个元素后，Dir调用路径和尾部的Clean。
//删除斜杠。
//如果`path`为空，则Dir返回“.”
//如果`path`为“.”，则Dir将该路径视为当前工作目录。
//如果`path`全部由分隔符组成，则Dir返回单个分隔符。
//除非是根目录，否则返回的路径不会以分隔符结尾。
func Dir(path string) string {
	if path == "." {
		return filepath.Dir(RealPath(path))
	}
	return filepath.Dir(path)
}

// IsEmpty 判断文件或者文件夹是否为空
func IsEmpty(path string) bool {
	stat, err := Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

// Ext 获取文件的扩展名，包含.开头
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// ExtName 获取文件扩展名，去除了"."
func ExtName(path string) string {
	return strings.TrimLeft(Ext(path), ".")
}

// TempDir 获取当前操作系统的临时目录
func TempDir(names ...string) string {
	path := tempDir
	for _, name := range names {
		path += Separator + name
	}
	return path
}
