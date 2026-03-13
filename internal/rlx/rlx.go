package rlx

import "sync/atomic"

type request struct {
	fn   func()
	done chan struct{}
}

var (
	reqCh        = make(chan request, 1024)
	onMainThread atomic.Bool
)

func Do(fn func()) {
	if onMainThread.Load() {
		fn()
		return
	}

	done := make(chan struct{})
	reqCh <- request{fn: fn, done: done}
	<-done
}

func Async(fn func()) {
	if onMainThread.Load() {
		fn()
		return
	}

	reqCh <- request{fn: fn}
}

func Call[T any](fn func() T) T {
	var out T
	Do(func() {
		out = fn()
	})
	return out
}

func EnterMainThread() {
	onMainThread.Store(true)
}

func ExitMainThread() {
	onMainThread.Store(false)
}

func Poll() {
	for {
		select {
		case req := <-reqCh:
			req.fn()
			if req.done != nil {
				close(req.done)
			}
		default:
			return
		}
	}
}
