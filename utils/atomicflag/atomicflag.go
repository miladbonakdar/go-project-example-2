package atomicflag

import "sync/atomic"

type AtomicFlag struct{ flag int32 }

func (b *AtomicFlag) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), i)
}

func (b *AtomicFlag) Get() bool {
	if atomic.LoadInt32(&(b.flag)) != 0 {
		return true
	}
	return false
}

func NewAtomicFlag() *AtomicFlag {
	flag := &AtomicFlag{}
	flag.Set(false)
	return flag
}

func NewAtomicFlagWithDefaultValue(defaultValue bool) *AtomicFlag {
	flag := &AtomicFlag{}
	flag.Set(defaultValue)
	return flag
}
