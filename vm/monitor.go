package vm

import "sync"

type (
	Monitor struct {
		m        *sync.Mutex
		entering []chan struct{}
		waiting  []chan struct{}
		owner    *Thread
		count    int
	}
)

func NewMonitor() *Monitor {
	return &Monitor{
		m:     &sync.Mutex{},
		owner: nil,
		count: 0,
	}
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

		entering := make(chan struct{})
		mon.entering = append(mon.entering, entering)
		mon.m.Unlock()

		<-entering
	}
}

func (mon *Monitor) Exit() {
	mon.m.Lock()
	defer mon.m.Unlock()

	mon.count--
	if mon.count == 0 {
		mon.owner = nil
		for _, w := range mon.entering {
			close(w)
		}
		mon.entering = nil
	}
}

func (mon *Monitor) Notify() {
	mon.m.Lock()
	defer mon.m.Unlock()

	if len(mon.waiting) == 0 {
		return
	}

	close(mon.waiting[0])
	mon.waiting = mon.waiting[1:]
}

func (mon *Monitor) NotifyAll() {
	mon.m.Lock()
	defer mon.m.Unlock()

	for _, n := range mon.waiting {
		close(n)
	}
	mon.waiting = nil
}

func (mon *Monitor) Wait() {
	mon.m.Lock()
	owner := mon.owner
	count := mon.count
	notify := make(chan struct{})
	mon.waiting = append(mon.waiting, notify)
	mon.m.Unlock()

	<-notify
	mon.Enter(owner, count)
}
