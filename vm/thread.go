package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
	"strings"
)

type Thread struct {
	vm          *VM
	java        *Instance
	derivedFrom *Thread
	frameStack  []*Frame

	result    interface{}
	unCatchEx *Instance
}

func NewThread(vm *VM) *Thread {
	return &Thread{vm: vm, java: nil}
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
		java:        thread.java,
		derivedFrom: thread,
		frameStack:  nil,
		unCatchEx:   nil,
	}
}

func (thread *Thread) Equal(t *Thread) bool {
	return t == thread || (thread.derivedFrom != nil && thread.derivedFrom.Equal(t))
}

func (thread *Thread) Execute(frame *Frame) (*Instance, error) {
	thread.frameStack = append(thread.frameStack, frame)
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
		trace = append(trace, fmt.Sprintf("%s.%s:%s", f.curClass.File().ThisClass(), *f.curMethod.Name(), *f.curMethod.Descriptor()))
	}

	return trace
}

func (thread *Thread) DumpFrameStack(showHeader bool) int {
	if showHeader {
		fmt.Println("------------ Frame stack ------------")
	}

	stackNum := 0
	if thread.derivedFrom != nil {
		stackNum += thread.derivedFrom.DumpFrameStack(false)
	}

	for _, f := range thread.frameStack {
		stackNum++
		indent := strings.Repeat("  ", stackNum)
		fmt.Printf(indent+"%s.%s:%s\n", f.curClass.File().ThisClass(), *f.curMethod.Name(), *f.curMethod.Descriptor())
	}

	if showHeader {
		fmt.Println("-------------------------------------")
	}

	return stackNum
}

func (thread *Thread) PushFrame(frame *Frame) {
	fmt.Printf("enter new frame: %s.%s:%s\n", frame.curClass.File().ThisClass(), *frame.curMethod.Name(), *frame.curMethod.Descriptor())

	if frame.CurrentClass().File().ThisClass() == "sun/reflect/Reflection" &&
		*(frame.CurrentMethod().Name()) == "isSameClassPackage" && len(frame.Locals()) >= 4 {
		c1 := frame.Locals()[1].(*Instance)
		c2 := frame.Locals()[3].(*Instance)

		fmt.Printf("--------------------------------- isSameClassPackage c=%s, m=%s\n", JavaStringToGoString(c1), JavaStringToGoString(c2))
	}

	thread.frameStack = append(thread.frameStack, frame)
}

func (thread *Thread) PopFrame() {
	fmt.Printf("leave frame: %s.%s:%s\n",
		thread.CurrentFrame().CurrentClass().File().ThisClass(),
		*thread.CurrentFrame().CurrentMethod().Name(),
		*thread.CurrentFrame().CurrentMethod().Descriptor(),
	)

	thread.frameStack = thread.frameStack[:len(thread.frameStack)-1]
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

func (thread *Thread) SetResult(value interface{}) {
	thread.result = value
}

func (thread *Thread) GetResult() interface{} {
	return thread.result
}

func (thread *Thread) HandleException(thrown *Instance) {
	for i := len(thread.frameStack) - 1; i >= 0; i-- {
		frame := thread.frameStack[i]
		handler := frame.FindCurrentExceptionHandler(thrown)

		if handler != nil {
			frame.JumpPC(*handler)
			frame.ClearOperand()
			frame.PushOperand(thrown)

			thread.frameStack = thread.frameStack[:i+1]
			return
		}
	}

	thread.unCatchEx = thrown
	thread.frameStack = nil
}
