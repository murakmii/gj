package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func ThrowableFillInStackTrace(thread *vm.Thread, args []interface{}) error {
	traces := thread.StackTrack()
	traceArray, traceSlice := vm.NewArray(thread.VM(), "[Ljava/lang/StackTraceElement;", len(traces))

	for i, t := range traces {
		traceSlice[i] = t.ToJava(thread.VM())
	}

	throwable := args[0].(*vm.Instance)

	traceName := "stackTrace"
	traceDesc := "[Ljava/lang/StackTraceElement;"
	throwable.PutField(&traceName, &traceDesc, traceArray)
	throwable.SetVMData(traces)

	thread.CurrentFrame().PushOperand(throwable)
	return nil
}

func ThrowableGetStackTraceDepth(thread *vm.Thread, args []interface{}) error {
	depth := int32(0)
	if traces, ok := args[0].(*vm.Instance).VMData().([]*vm.StackTraceElement); ok {
		depth = int32(len(traces))
	}

	thread.CurrentFrame().PushOperand(depth)
	return nil
}

func ThrowableGetStackTraceElement(thread *vm.Thread, args []interface{}) error {
	throwable := args[0].(*vm.Instance)
	traces, ok := throwable.VMData().([]*vm.StackTraceElement)
	if !ok {
		return fmt.Errorf("%s has NOT stack trace", throwable.Class().File().ThisClass())
	}

	thread.CurrentFrame().PushOperand(traces[args[1].(int32)].ToJava(thread.VM()))
	return nil
}
