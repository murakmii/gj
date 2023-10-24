package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"sync"
)

type (
	Thread struct {
		vm          *VM
		name        string
		main        bool
		daemon      bool
		java        *Instance
		derivedFrom *Thread
		frameStack  []*Frame
		syncStack   []*Instance
		alive       bool
		unCatchEx   *Instance

		interLock    *sync.Mutex
		interrupted  bool
		interWatcher []chan struct{}
	}

	ThreadResult struct {
		Thread    *Thread
		Err       error
		UnCatchEx *Instance
	}

	ThreadExecutor struct {
		lock         *sync.Mutex
		executingNum int
		daemonNum    int
		result       chan *ThreadResult
	}
)

func NewThread(vm *VM, name string, main, daemon bool) *Thread {
	return &Thread{
		vm:        vm,
		name:      name,
		main:      main,
		daemon:    daemon,
		alive:     true,
		interLock: &sync.Mutex{},
	}
}

func (thread *Thread) JavaThread() *Instance {
	return thread.java
}

func (thread *Thread) SetJavaThread(java *Instance) {
	thread.java = java
}

func (thread *Thread) VM() *VM {
	return thread.vm
}

func (thread *Thread) Derive() *Thread {
	return &Thread{
		vm:          thread.vm,
		name:        thread.name,
		java:        thread.java,
		derivedFrom: thread,
	}
}

func (thread *Thread) Name() string {
	return thread.name
}

func (thread *Thread) SetName(name string) {
	thread.name = name
}

func (thread *Thread) IsAlive() bool {
	return thread.alive
}

func (thread *Thread) IsDaemon() bool {
	return thread.daemon
}

func (thread *Thread) Equal(t *Thread) bool {
	return t == thread || (thread.derivedFrom != nil && thread.derivedFrom.Equal(t))
}

func (thread *Thread) Execute(frame *Frame) (*Instance, error) {
	thread.PushFrame(frame)

	for len(thread.frameStack) > 0 {
		curFrame := thread.frameStack[len(thread.frameStack)-1]

		if err := ExecInstr(thread, curFrame, curFrame.NextInstr()); err != nil {
			return nil, err
		}
	}

	return thread.unCatchEx, nil
}

func (thread *Thread) ExecMethod(class *Class, method *class_file.MethodInfo) error {
	curFrame := thread.CurrentFrame()
	args := curFrame.PopOperands(method.NumArgs())

	if method.IsNative() {
		return CallNativeMethod(thread, class, method, args)
	}

	thread.PushFrame(NewFrame(class, method).SetLocals(args))
	return nil
}

func (thread *Thread) StackTrack() []string {
	var trace []string

	if thread.derivedFrom != nil {
		for _, t := range thread.derivedFrom.StackTrack() {
			trace = append(trace, t)
		}
	}

	for _, f := range thread.frameStack {
		trace = append(trace, fmt.Sprintf("%s.%s:%s", f.curClass.File().ThisClass(), *f.curMethod.Name(), f.curMethod.Descriptor()))
	}

	return trace
}

func (thread *Thread) PushFrame(frame *Frame) {
	var syncObj *Instance

	if frame.CurrentMethod().IsSync() {
		if frame.CurrentMethod().IsStatic() {
			syncObj = frame.CurrentClass().Java()
		} else {
			syncObj = frame.Locals()[0].(*Instance)
		}
		syncObj.Monitor().Enter(thread, -1)
	}

	thread.frameStack = append(thread.frameStack, frame)
	thread.syncStack = append(thread.syncStack, syncObj)
}

func (thread *Thread) PopFrame() {
	idx := len(thread.frameStack) - 1

	if thread.frameStack[idx].CurrentMethod().IsSync() {
		thread.syncStack[idx].Monitor().Exit(thread)
	}

	thread.frameStack = thread.frameStack[:idx]
	thread.syncStack = thread.syncStack[:idx]
}

func (thread *Thread) InvokerFrame() *Frame {
	if len(thread.frameStack) < 2 {
		return nil
	}
	return thread.frameStack[len(thread.frameStack)-2]
}

func (thread *Thread) CurrentFrame() *Frame {
	if len(thread.frameStack) == 0 {
		return nil
	}
	return thread.frameStack[len(thread.frameStack)-1]
}

func (thread *Thread) HandleException(thrown *Instance) {
	for len(thread.frameStack) > 0 {
		frame := thread.CurrentFrame()
		handler := frame.FindCurrentExceptionHandler(thrown)

		if handler != nil {
			frame.JumpPC(*handler)
			frame.ClearOperand()
			frame.PushOperand(thrown)
			return
		}

		thread.PopFrame()
	}

	thread.unCatchEx = thrown
}

func (thread *Thread) Interrupt() {
	thread.interLock.Lock()
	defer thread.interLock.Unlock()

	for _, w := range thread.interWatcher {
		close(w)
	}

	thread.interrupted = len(thread.interWatcher) > 0
	thread.interWatcher = nil
}

func (thread *Thread) WatchInterruption() <-chan struct{} {
	thread.interLock.Lock()
	defer thread.interLock.Unlock()

	watcher := make(chan struct{})
	thread.interWatcher = append(thread.interWatcher, watcher)

	return watcher
}

func (thread *Thread) UnWatchInterruption(watcher <-chan struct{}) {
	thread.interLock.Lock()
	defer thread.interLock.Unlock()

	for i, w := range thread.interWatcher {
		if w != watcher {
			continue
		}
		thread.interWatcher = append(thread.interWatcher[:i], thread.interWatcher[i+1:]...)
		break
	}
}

func NewThreadExecutor() *ThreadExecutor {
	return &ThreadExecutor{lock: &sync.Mutex{}, result: make(chan *ThreadResult)}
}

// Start goroutine to execute 'frame' on 'thread'
func (executor *ThreadExecutor) Start(thread *Thread, frame *Frame) {
	executor.lock.Lock()
	defer executor.lock.Unlock()

	executor.executingNum++
	if thread.daemon {
		executor.daemonNum++
	}

	go func() {
		unCatch, err := thread.Execute(frame)
		thread.alive = false

		thread.JavaThread().Monitor().Enter(thread, -1)
		thread.JavaThread().Monitor().NotifyAll(thread)
		thread.JavaThread().Monitor().Exit(thread)

		executor.lock.Lock()
		executor.executingNum--
		if thread.IsDaemon() {
			executor.daemonNum--
		}
		done := executor.executingNum-executor.daemonNum == 0
		executor.lock.Unlock()

		executor.result <- &ThreadResult{
			Thread:    thread,
			Err:       err,
			UnCatchEx: unCatch,
		}

		if done {
			close(executor.result)
		}
	}()
}

// Receiving result of each thread execution.
// If all non-daemon threads finished, channel will be closed.
func (executor *ThreadExecutor) Wait() <-chan *ThreadResult {
	return executor.result
}
