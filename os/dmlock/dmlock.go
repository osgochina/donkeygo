package dmlock

var (
	// 默认锁 locker.
	locker = New()
)

// Lock 对key加锁，如果key已经存在锁，则会阻塞等待锁
func Lock(key string) {
	locker.Lock(key)
}

// TryLock 尝试对key加锁，如果加锁成功返回true，加锁失败返回false
func TryLock(key string) bool {
	return locker.TryLock(key)
}

// Unlock 对key解锁
func Unlock(key string) {
	locker.Unlock(key)
}

// RLock 都key加读锁，如果key已存在写锁，则会阻塞等待
func RLock(key string) {
	locker.RLock(key)
}

// TryRLock 尝试对key加读锁，如果加锁成功返回true，失败返回false
func TryRLock(key string) bool {
	return locker.TryRLock(key)
}

// RUnlock 对key进行读解锁
func RUnlock(key string) {
	locker.RUnlock(key)
}

// LockFunc 对key加锁，并且执行方法f，执行成功自动解锁
func LockFunc(key string, f func()) {
	locker.LockFunc(key, f)
}

// RLockFunc 对key加读锁执行方法f，执行成功自动解锁
func RLockFunc(key string, f func()) {
	locker.RLockFunc(key, f)
}

// TryLockFunc 尝试对key加锁，如果加锁成功则执行方法f，执行成功返回true，执行失败返回false
func TryLockFunc(key string, f func()) bool {
	return locker.TryLockFunc(key, f)
}

// TryRLockFunc 尝试对key加读锁，如果加锁成功则执行方法f，执行成功返回true，加锁失败或执行失败都返回false
func TryRLockFunc(key string, f func()) bool {
	return locker.TryRLockFunc(key, f)
}

// Remove 移除在key身上的锁
func Remove(key string) {
	locker.Remove(key)
}
