package utils

import "sync/atomic"

type AtomicBool int32

func (a *AtomicBool) Set(value bool) (swapped bool) {
	if value {
		return atomic.SwapInt32((*int32)(a), 1) == 0
	}
	return atomic.SwapInt32((*int32)(a), 0) == 1
}

func (a *AtomicBool) Get() bool {
	return atomic.LoadInt32((*int32)(a)) != 0
}
