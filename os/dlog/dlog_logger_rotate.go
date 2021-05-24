package dlog

import (
	"fmt"
	"github.com/gogf/gf/os/gmlock"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/encoding/dcompress"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/text/dregex"
	"time"
)

//备份日志文件
func (that *Logger) rotateFileBySize(now time.Time) {
	if that.config.RotateSize <= 0 {
		return
	}
	if err := that.doRotateFile(that.getFilePath(now)); err != nil {
		// panic(err)
		intlog.Error(err)
	}
}

//备份日志文件
func (that *Logger) doRotateFile(filePath string) error {
	memoryLockKey := "dlog.doRotateFile:" + filePath
	if !gmlock.TryLock(memoryLockKey) {
		return nil
	}
	defer gmlock.Unlock(memoryLockKey)

	//最大备份数，如果为0，则表示不备份，把该文件删除
	if that.config.RotateBackupLimit == 0 {
		if err := dfile.Remove(filePath); err != nil {
			return err
		}
		intlog.Printf(`%d size exceeds, no backups set, remove original logging file: %s`, that.config.RotateSize, filePath)
		return nil
	}
	// 为备份做准备，留下原始原件的信息
	var (
		dirPath     = dfile.Dir(filePath)
		fileName    = dfile.Name(filePath)
		fileExtName = dfile.ExtName(filePath)
		newFilePath = ""
	)
	for {
		var (
			now   = dtime.Now()
			micro = now.Microsecond() % 1000
		)
		if micro == 0 {
			micro = 101
		} else {
			for micro < 100 {
				micro *= 10
			}
		}
		//新文件名
		newFilePath = dfile.Join(
			dirPath,
			fmt.Sprintf(
				`%s.%s%d.%s`,
				fileName, now.Format("YmdHisu"), micro, fileExtName,
			),
		)
		if !dfile.Exists(newFilePath) {
			break
		} else {
			intlog.Printf(`rotation file exists, continue: %s`, newFilePath)
		}
	}
	//把老的文件改名
	if err := dfile.Rename(filePath, newFilePath); err != nil {
		return err
	}
	return nil
}

// 通过定时器，检查日志文件的备份及压缩情况
func (that *Logger) rotateChecksTimely() {
	//每次执行完毕，重新把定时方法加入
	defer dtimer.AddOnce(that.config.RotateCheckInterval, that.rotateChecksTimely)

	//旋转文件备份未启动
	if that.config.RotateSize <= 0 && that.config.RotateExpire == 0 {
		intlog.Printf(
			"logging rotation ignore checks: RotateSize: %d, RotateExpire: %s",
			that.config.RotateSize, that.config.RotateExpire.String(),
		)
		return
	}
	//加入内存锁
	memoryLockKey := "dlog.rotateChecksTimely:" + that.config.Path
	if !gmlock.TryLock(memoryLockKey) {
		return
	}
	defer gmlock.Unlock(memoryLockKey)

	var (
		now      = time.Now()
		pattern  = "*.log, *.gz"
		files, _ = dfile.ScanDirFile(that.config.Path, pattern, true)
	)
	intlog.Printf("logging rotation start checks: %+v", files)

	// =============================================================
	// 检查旋转文件是否过期
	// =============================================================
	if that.config.RotateExpire > 0 {
		var (
			mtime         time.Time
			subDuration   time.Duration
			expireRotated bool
		)
		//如果文件已过期，则移动文件
		for _, file := range files {
			if dfile.ExtName(file) == "gz" {
				continue
			}
			mtime = dfile.MTime(file)
			subDuration = now.Sub(mtime)
			if subDuration > that.config.RotateExpire {
				expireRotated = true
				intlog.Printf(
					`%v - %v = %v > %v, rotation expire logging file: %s`,
					now, mtime, subDuration, that.config.RotateExpire, file,
				)
				if err := that.doRotateFile(file); err != nil {
					intlog.Error(err)
				}
			}
		}
		//移动成功后重新扫描该文件夹下文件列表
		if expireRotated {
			// Update the files array.
			files, _ = dfile.ScanDirFile(that.config.Path, pattern, true)
		}
	}

	// =============================================================
	// 文件压缩
	// =============================================================
	needCompressFileArray := darray.NewStrArray()
	if that.config.RotateBackupCompress > 0 {
		for _, file := range files {
			// Eg: access.20200326101301899002.log.gz
			if dfile.ExtName(file) == "gz" {
				continue
			}
			// Eg:
			// access.20200326101301899002.log
			if dregex.IsMatchString(`.+\.\d{20}\.log`, dfile.Basename(file)) {
				needCompressFileArray.Append(file)
			}
		}
		if needCompressFileArray.Len() > 0 {
			needCompressFileArray.Iterator(func(_ int, path string) bool {
				err := dcompress.GzipFile(path, path+".gz")
				if err == nil {
					intlog.Printf(`compressed done, remove original logging file: %s`, path)
					if err = dfile.Remove(path); err != nil {
						intlog.Print(err)
					}
				} else {
					intlog.Print(err)
				}
				return true
			})
			// Update the files array.
			files, _ = dfile.ScanDirFile(that.config.Path, pattern, true)
		}
	}

	// =============================================================
	// 备份计数限制和过期检查。
	// =============================================================
	var (
		backupFilesMap          = make(map[string]*darray.SortedArray)
		originalLoggingFilePath = ""
	)
	if that.config.RotateBackupLimit > 0 || that.config.RotateBackupExpire > 0 {
		for _, file := range files {
			originalLoggingFilePath, _ = dregex.ReplaceString(`\.\d{20}`, "", file)
			if backupFilesMap[originalLoggingFilePath] == nil {
				backupFilesMap[originalLoggingFilePath] = darray.NewSortedArray(func(a, b interface{}) int {
					// Sorted by rotated/backup file mtime.
					// The old rotated/backup file is put in the head of array.
					file1 := a.(string)
					file2 := b.(string)
					result := dfile.MTimestampMilli(file1) - dfile.MTimestampMilli(file2)
					if result <= 0 {
						return -1
					}
					return 1
				})
			}
			// Check if this file a rotated/backup file.
			if dregex.IsMatchString(`.+\.\d{20}\.log`, dfile.Basename(file)) {
				backupFilesMap[originalLoggingFilePath].Add(file)
			}
		}
		intlog.Printf(`calculated backup files map: %+v`, backupFilesMap)
		for _, array := range backupFilesMap {
			diff := array.Len() - that.config.RotateBackupLimit
			for i := 0; i < diff; i++ {
				path, _ := array.PopLeft()
				intlog.Printf(`remove exceeded backup limit file: %s`, path)
				if err := dfile.Remove(path.(string)); err != nil {
					intlog.Print(err)
				}
			}
		}
		// 备份过期文件
		if that.config.RotateBackupExpire > 0 {
			var (
				mtime       time.Time
				subDuration time.Duration
			)
			for _, array := range backupFilesMap {
				array.Iterator(func(_ int, v interface{}) bool {
					path := v.(string)
					mtime = dfile.MTime(path)
					subDuration = now.Sub(mtime)
					if subDuration > that.config.RotateBackupExpire {
						intlog.Printf(
							`%v - %v = %v > %v, remove expired backup file: %s`,
							now, mtime, subDuration, that.config.RotateBackupExpire, path,
						)
						if err := dfile.Remove(path); err != nil {
							intlog.Print(err)
						}
						return true
					} else {
						return false
					}
				})
			}
		}
	}
}
