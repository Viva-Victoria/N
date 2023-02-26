package sync

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNWaitGroup_NewWaitGroup(t *testing.T) {
	nwg := NewWaitGroup(0)
	assert.Equal(t, int64(0), nwg.total.Load())
	assertNoLock(t, nwg.Wait, time.Millisecond)
	assertNoLock(t, func() {
		require.NoError(t, nwg.WaitContext(context.Background()))
	}, time.Millisecond)
	assertNoLock(t, func() {
		require.NoError(t, nwg.WaitTimeout(time.Millisecond))
	}, time.Millisecond)
}

func TestNWaitGroup_Add(t *testing.T) {
	t.Run("simple+closed", func(t *testing.T) {
		nwg := NewWaitGroup(0)
		require.NoError(t, nwg.Add(3))

		nwg.Done(2)
		require.NoError(t, nwg.Add(2))

		nwg.Done(3)
		require.Error(t, ErrWaitGroupClosed, nwg.Add(2))
	})
	t.Run("negative", func(t *testing.T) {
		nwg := NewWaitGroup(0)
		require.NoError(t, nwg.Add(3))

		require.NoError(t, nwg.Add(-2))
		require.NoError(t, nwg.Add(2))

		require.NoError(t, nwg.Add(-3))
		require.Error(t, ErrWaitGroupClosed, nwg.Add(2))
	})
}

func TestNWaitGroup_Done(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		nwg := NewWaitGroup(0)
		require.Panics(t, func() {
			nwg.Done(1)
		})
	})
	t.Run("not-panic", func(t *testing.T) {
		nwg := NewWaitGroup(1)
		require.NotPanics(t, func() {
			nwg.Done(1)
		})
	})
	t.Run("simple+closed", func(t *testing.T) {
		nwg := NewWaitGroup(0)
		_ = nwg.Add(3)
		nwg.Done(2)

		_ = nwg.Add(2)
		nwg.Done(3)

		assert.Equal(t, true, nwg.closed.Load())
		assertNoLock(t, nwg.Wait, time.Millisecond)

		assert.NotPanics(t, func() {
			nwg.Done(1)
		})
	})
}

func TestNWaitGroup_Wait(t *testing.T) {
	t.Run("no-lock", func(t *testing.T) {
		nwg := NewWaitGroup(0)
		_ = nwg.Add(1)
		go func() {
			<-time.After(time.Millisecond * 47)
			nwg.Done(1)
		}()

		assertNoLock(t, nwg.Wait, time.Millisecond*50)
	})
	t.Run("lock", func(t *testing.T) {
		nwg := NewWaitGroup(1)
		assertLock(t, nwg.Wait, time.Millisecond)
	})
}

func TestNWaitGroup_WaitContext(t *testing.T) {
	t.Run("no-lock", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			nwg := NewWaitGroup(0)
			_ = nwg.Add(1)
			go func() {
				<-time.After(time.Millisecond * 47)
				nwg.Done(1)
			}()

			assertNoLock(t, func() {
				require.NoError(t, nwg.WaitContext(context.Background()))
			}, time.Millisecond*50)
		})
		t.Run("timeout-ctx", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertNoLock(t, func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
				defer cancel()

				_ = nwg.WaitContext(ctx)
			}, time.Millisecond*6)
		})
	})
	t.Run("lock", func(t *testing.T) {
		t.Run("background", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertLock(t, func() {
				_ = nwg.WaitContext(context.Background())
			}, time.Millisecond)
		})
		t.Run("timeout-ctx", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertLock(t, func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
				defer cancel()

				_ = nwg.WaitContext(ctx)
			}, time.Millisecond*45)
		})
	})
}

func TestNWaitGroup_WaitTimeout(t *testing.T) {
	t.Run("no-lock", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			nwg := NewWaitGroup(0)
			_ = nwg.Add(1)
			go func() {
				<-time.After(time.Millisecond * 47)
				nwg.Done(1)
			}()

			assertNoLock(t, func() {
				require.NoError(t, nwg.WaitTimeout(time.Millisecond*49))
			}, time.Millisecond*51)
		})
		t.Run("exceed", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertNoLock(t, func() {
				require.Error(t, ErrTimeoutExceed, nwg.WaitTimeout(time.Millisecond*5))
			}, time.Millisecond*6)
		})
	})
	t.Run("lock", func(t *testing.T) {
		t.Run("big-limit", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertLock(t, func() {
				_ = nwg.WaitTimeout(time.Millisecond * 50)
			}, time.Millisecond)
		})
		t.Run("good-limit", func(t *testing.T) {
			nwg := NewWaitGroup(1)
			assertLock(t, func() {
				require.Error(t, ErrTimeoutExceed, nwg.WaitTimeout(time.Millisecond*40))
			}, time.Millisecond*40)
		})
	})
}

func assertNoLock(t *testing.T, f func(), timeout time.Duration) {
	start := time.Now()
	returnChan := make(chan struct{})
	go func() {
		defer close(returnChan)
		f()
	}()

	select {
	case <-returnChan:
		t.Logf("executed with %v (limit %v)", time.Since(start), timeout)
		return
	case <-time.After(timeout):
		t.Fatalf("should be executed with %v, but still locked or executing", timeout)
		return
	}
}

func assertLock(t *testing.T, f func(), timeout time.Duration) {
	start := time.Now()
	returnChan := make(chan struct{})
	go func() {
		defer close(returnChan)
		f()
	}()

	select {
	case <-returnChan:
		t.Fatalf("should be locked for %v, but executed by %v", timeout, time.Since(start))
		return
	case <-time.After(timeout):
		t.Logf("limit %v exceed, f() still locked", timeout)
		return
	}
}
