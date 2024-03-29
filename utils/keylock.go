package utils

import (
	"sync"
	"sync/atomic"
)

var mutexPool = sync.Pool{
	New: func() interface{} {
		return &sync.Mutex{}
	},
}

type LockObj struct {
	Lock *sync.Mutex
	Num  int64
}

type KeyLock struct {
	globalLock sync.Mutex
	locks      map[interface{}]*LockObj
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks: make(map[interface{}]*LockObj),
	}
}

func (l *KeyLock) getLock(key interface{}) *sync.Mutex {

	l.globalLock.Lock()
	defer l.globalLock.Unlock()

	if lockObj, ok := l.locks[key]; ok {
		// 紀錄waiting中的數量
		atomic.AddInt64(&lockObj.Num, 1)
		return lockObj.Lock
	}
	lock := mutexPool.Get().(*sync.Mutex)
	l.locks[key] = &LockObj{
		Lock: lock,
		Num:  1,
	}
	return lock
}

func (l *KeyLock) Lock(key interface{}) {
	l.getLock(key).Lock()
}

func (l *KeyLock) Unlock(key interface{}) {
	l.globalLock.Lock()
	defer l.globalLock.Unlock()

	l.locks[key].Lock.Unlock()
	atomic.AddInt64(&l.locks[key].Num, -1)
	//clean
	for _, v := range l.locks {
		// 判斷沒有等待中的鎖 釋放緩存
		if v.Num <= 0 {
			// 放回池中, 供其他協程使用
			mutexPool.Put(l.locks[key].Lock)
			delete(l.locks, key)
		}
	}
}
