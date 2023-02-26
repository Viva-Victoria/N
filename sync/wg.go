package sync

import (
	"context"
	"errors"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	ErrTimeoutExceed   = errors.New("timeout exceed")
	ErrWaitGroupClosed = errors.New("WaitGroup closed")
)

type WaitGroup interface {
	Add(delta int) error
	Done(delta int)
	Wait()
	WaitContext(ctx context.Context) error
	WaitTimeout(timeout time.Duration) error
}

type NWaitGroup struct {
	closed   atomic.Bool
	total    atomic.Int64
	wg       sync.WaitGroup
	waitChan chan struct{}
}

func NewWaitGroup(initial int) *NWaitGroup {
	nwg := &NWaitGroup{
		closed:   *atomic.NewBool(false),
		total:    *atomic.NewInt64(0),
		wg:       sync.WaitGroup{},
		waitChan: make(chan struct{}),
	}
	_ = nwg.Add(initial)
	return nwg
}

func (N *NWaitGroup) Add(delta int) error {
	if N.closed.Load() {
		return ErrWaitGroupClosed
	}

	N.wg.Add(delta)
	N.total.Add(int64(delta))
	return nil
}

func (N *NWaitGroup) Done(delta int) {
	if N.closed.Load() {
		return
	}

	for i := 0; i < delta; i++ {
		N.wg.Done()
	}

	if N.total.Sub(int64(delta)) == 0 {
		N.close()
	}
}

func (N *NWaitGroup) Wait() {
	N.wg.Wait()
}

func (N *NWaitGroup) WaitContext(ctx context.Context) error {
	if N.total.Load() == 0 {
		return nil
	}

	select {
	case <-N.waitChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (N *NWaitGroup) WaitTimeout(timeout time.Duration) error {
	if N.total.Load() == 0 {
		return nil
	}

	select {
	case <-N.waitChan:
		return nil
	case <-time.After(timeout):
		return ErrTimeoutExceed
	}
}

func (N *NWaitGroup) close() {
	N.closed.Store(true)
	close(N.waitChan)
}
