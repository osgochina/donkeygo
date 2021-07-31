package dmutex

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"math"
	"runtime"
)

// Mutex 高级互斥锁
// 使用atomic + channel实现的高级互斥锁模块，支持更丰富的互斥锁特性
// 互斥锁对象支持读写控制，互斥锁功能逻辑与标准库sync.RWMutex类似，可并发读但不可并发写
type Mutex struct {
	state   *dtype.Int32  //锁的状态，-1 表示写锁中，>1 表示读锁中
	writer  *dtype.Int32  // 当前等待写锁的个数
	reader  *dtype.Int32  // 当前等待读锁的个数
	writing chan struct{} // 写锁使用的通道
	reading chan struct{} // 读锁使用的通道
}

// New 新建一个高级互斥锁
func New() *Mutex {
	return &Mutex{
		state:   dtype.NewInt32(),
		writer:  dtype.NewInt32(),
		reader:  dtype.NewInt32(),
		writing: make(chan struct{}, 1),
		reading: make(chan struct{}, math.MaxInt32),
	}
}

// Lock 加写锁，如果锁已经被其他的goroutine占用，则会阻塞等待锁可用
func (that *Mutex) Lock() {
	for {
		//加锁成功
		if that.state.Cas(0, -1) {
			return
		}
		// 加锁失败，等待下一次加锁的机会,
		that.writer.Add(1)
		<-that.writing
	}
}

// Unlock 解锁高级互斥锁上的写锁，支持重复调用
func (that *Mutex) Unlock() {
	//解锁，如果当前状态不是被写锁阻塞，则说明都不做,
	if that.state.Cas(-1, 0) {
		//有多个协程会进度到这个逻辑
		var n int32
		for {
			// 如果写锁解锁，则优先检查等待读锁的协程，如果有读锁在阻塞等待锁，则使用抢占式的来激活它们，让他们获取到锁
			// 有多个协程在等待读锁，则会依次激活他们，让他们获取到锁
			if n = that.reader.Val(); n > 0 {
				if that.reader.Cas(n, 0) {
					for ; n > 0; n-- {
						that.reading <- struct{}{}
					}
					break
				} else {
					// 如果
					runtime.Gosched()
				}
			} else {
				break
			}
		}

		// 如果有协程在等待写锁，给与该等待着一个获取锁的机会
		if n = that.writer.Val(); n > 0 {
			if that.writer.Cas(n, n-1) {
				that.writing <- struct{}{}
			}
		}
	}
}

// TryLock 非阻塞尝试获取写锁，获取锁成功返回true，获取失败返回false
func (that *Mutex) TryLock() bool {
	if that.state.Cas(0, -1) {
		return true
	}
	return false
}

// RLock 给高级互斥锁加上读锁，如果已存在写锁，则需要等待写锁完成在执行
func (that *Mutex) RLock() {
	var n int32
	for {
		// 因为读锁是>0,0表示没有锁，所以这里的逻辑是，检查是否有写锁
		if n = that.state.Val(); n >= 0 {
			// 检查读锁，并加锁,如果读锁加锁失败，则让出cpu，等待下一次调度重新执行
			if that.state.Cas(n, n+1) {
				return
			} else {
				runtime.Gosched()
			}
		} else {
			//如果有写锁，则需要等待写锁解锁后，再抢占读锁
			that.reader.Add(1)
			<-that.reading
		}
	}
}

// RUnlock 解锁高级互斥锁上的读锁，可以重复调用
func (that *Mutex) RUnlock() {
	var n int32
	for {
		//获取当前是否有读锁，如果有读锁，则解锁，如果解锁不成功，则让出cpu，等待下次调度再次解锁
		// 如果当前不是读锁，而是写锁，则跳过
		if n = that.state.Val(); n >= 1 {
			if that.state.Cas(n, n-1) {
				break
			} else {
				runtime.Gosched()
			}
		} else {
			break
		}
	}

	// n == 1 表示高级互斥锁已经解锁，状态为0
	// 检查当前阻塞等待的写锁的协程，如果有阻塞等待的协程，则主动通知该协程可以去获取写锁了
	if n == 1 {
		if n = that.writer.Val(); n > 0 {
			if that.writer.Cas(n, n-1) {
				that.writing <- struct{}{}
			}
		}
	}
}

// TryRLock 尝试给高级互斥锁增加读锁，如果添加成功，则返回true，加锁失败返回false
func (that *Mutex) TryRLock() bool {
	var n int32
	for {
		// 判断当前锁是不是在可以加读锁的状态
		if n = that.state.Val(); n >= 0 {
			//尝试加锁
			if that.state.Cas(n, n+1) {
				return true
			} else {
				runtime.Gosched()
			}
		} else {
			return false
		}
	}
}

// IsLocked 判断当前锁是否处于加锁状态
func (that *Mutex) IsLocked() bool {
	return that.state.Val() != 0
}

// IsWLocked 判断当前锁是否是处于写锁状态
func (that *Mutex) IsWLocked() bool {
	return that.state.Val() < 0
}

// IsRLocked 判断当前锁是否处于读锁状态
func (that *Mutex) IsRLocked() bool {
	return that.state.Val() > 0
}

// LockFunc 使用写锁执行一段代码，不用担心panic，因为使用defer解锁，所以无论如何都会解锁成功
func (that *Mutex) LockFunc(f func()) {
	that.Lock()
	defer that.Unlock()
	f()
}

// RLockFunc 使用读锁直营一段代码，一定会解锁成功
func (that *Mutex) RLockFunc(f func()) {
	that.RLock()
	defer that.RUnlock()
	f()
}

// TryLockFunc 尝试加锁并执行方法，加锁成功则会执行方法，加锁失败什么都不做，
// 返回值result表示加锁是否成功
func (that *Mutex) TryLockFunc(f func()) (result bool) {
	if that.TryLock() {
		result = true
		defer that.Unlock()
		f()
	}
	return
}

// TryRLockFunc 尝试加读锁执行一段方法，加锁成功则会执行，加锁失败则什么都不做，返回值result表示加锁是否成功
func (that *Mutex) TryRLockFunc(f func()) (result bool) {
	if that.TryRLock() {
		result = true
		defer that.RUnlock()
		f()
	}
	return
}
