package vm

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	// Implementation for synchronize, wait, notify and notifyAll
	Monitor struct {
		m        *sync.Mutex
		entering []chan struct{}
		waiting  []chan struct{}
		owner    *Thread
		count    int
	}
)

func NewMonitor() *Monitor {
	return &Monitor{m: &sync.Mutex{}}
}

func (mon *Monitor) Enter(thread *Thread, count int) {
	for {
		mon.m.Lock()

		if mon.owner == nil || mon.owner == thread {
			mon.owner = thread
			if count == -1 {
				mon.count++
			} else {
				mon.count = count
			}

			mon.m.Unlock()
			return
		}

		// If thread does NOT become owner, register channel to wait releasing.
		// This channel will be closed in 'Exit' method by owner.
		entering := make(chan struct{})
		mon.entering = append(mon.entering, entering)
		mon.m.Unlock()

		<-entering // Owner released monitor. Try to acquire ownership in next loop.
	}
}

func (mon *Monitor) Exit(owner *Thread) error {
	mon.m.Lock()
	defer mon.m.Unlock()

	if err := mon.assertOwner(owner); err != nil {
		return err
	}

	mon.count--
	if mon.count == 0 {
		mon.owner = nil
		for _, w := range mon.entering {
			close(w) // Resume threads that is waiting releasing.
		}
		mon.entering = nil
	}

	return nil
}

func (mon *Monitor) Wait(owner *Thread, timeoutMs int) (bool, error) {
	mon.m.Lock()

	if err := mon.assertOwner(owner); err != nil {
		return false, err
	}

	// Save count to restore it after re-entering.
	count := mon.count
	mon.count = 0
	mon.owner = nil

	// 'notify' will be closed by 'Notify' or 'NotifyAll' methods
	notify := make(chan struct{})
	mon.waiting = append(mon.waiting, notify)

	mon.m.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	if timeoutMs > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
	}
	defer cancel()

	var interrupted bool
	inter := owner.WatchInterruption()
	defer owner.UnWatchInterruption(inter)

	select {
	case <-notify:
	case <-inter:
		interrupted = true
	case <-ctx.Done():
		mon.m.Lock()
		for i, n := range mon.waiting {
			if n != notify {
				continue
			}
			mon.waiting = append(mon.waiting[:i], mon.waiting[i+1:]...)
			break
		}
		mon.m.Unlock()
	}

	mon.Enter(owner, count)
	return interrupted, nil
}

func (mon *Monitor) Notify(owner *Thread) error {
	mon.m.Lock()
	defer mon.m.Unlock()

	if err := mon.assertOwner(owner); err != nil {
		return err
	}

	if len(mon.waiting) == 0 {
		return nil
	}

	close(mon.waiting[0])
	mon.waiting = mon.waiting[1:]

	return nil
}

func (mon *Monitor) NotifyAll(owner *Thread) error {
	mon.m.Lock()
	defer mon.m.Unlock()

	if err := mon.assertOwner(owner); err != nil {
		return err
	}

	for _, n := range mon.waiting {
		close(n)
	}
	mon.waiting = nil

	return nil
}

func (mon *Monitor) assertOwner(owner *Thread) error {
	if mon.owner != owner {
		fmt.Printf("try ownership: %p, owner: %p\n", owner, mon.owner)
		return fmt.Errorf("try ownership = %s(%p), is NOT %s(%p)", owner.name, owner, mon.owner.name, mon.owner)
	}
	return nil
}
