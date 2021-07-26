package dfsnotify

//返回事件的字符串说明
func (that *Event) String() string {
	return that.event.String()
}

// IsCreate 是否是文件创建事件
func (that *Event) IsCreate() bool {
	return that.Op == 1 || that.Op&CREATE == CREATE
}

// IsWrite 是否是文件写入事件
func (that *Event) IsWrite() bool {
	return that.Op&WRITE == WRITE
}

// IsRemove 是否是文件删除事件
func (that *Event) IsRemove() bool {
	return that.Op&REMOVE == REMOVE
}

// IsRename 是否是文件改名事件
func (that *Event) IsRename() bool {
	return that.Op&RENAME == RENAME
}

// IsChmod 是否是文件权限修改事件
func (that *Event) IsChmod() bool {
	return that.Op&CHMOD == CHMOD
}
