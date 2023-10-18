package vm

import "sync"

type Monitor struct {
	cond  *sync.Cond
	owner *Thread
	count int
}

func NewMonitor() *Monitor {
	return &Monitor{
		cond:  sync.NewCond(&sync.Mutex{}),
		owner: nil,
		count: 0,
	}
}

func (monitor *Monitor) Enter(thread *Thread) {
	for {
		monitor.cond.L.Lock()

		if monitor.count == 0 || monitor.owner == thread {
			monitor.owner = thread
			monitor.count++
			monitor.cond.L.Unlock()
			return
		}

		monitor.cond.Wait()
	}
}

func (monitor *Monitor) Exit() {
	monitor.cond.L.Lock()
	defer monitor.cond.L.Unlock()

	monitor.count--
	if monitor.count == 0 {
		monitor.owner = nil
		monitor.cond.Broadcast()
	}
}
