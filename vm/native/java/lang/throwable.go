package lang

import (
	"fmt"
	"github.com/murakmii/gojiai/vm"
)

func init() {
	class := "java/lang/Throwable"

	vm.NativeMethods.Register(class, "fillInStackTrace", "(I)Ljava/lang/Throwable;", func(thread *vm.Thread, args []interface{}) error {
		traces := thread.StackTrack()
		traceArray, traceSlice := vm.NewArray(thread.VM(), "[Ljava/lang/StackTraceElement;", len(traces))

		for i, t := range traces {
			traceSlice[i] = t.ToJava(thread.VM())
		}

		throwable := args[0].(*vm.Instance)
		throwable.PutField("stackTrace", "[Ljava/lang/StackTraceElement;", traceArray)
		throwable.ToBeThrowable(traces)

		thread.CurrentFrame().PushOperand(throwable)
		return nil
	})

	vm.NativeMethods.Register(class, "getStackTraceDepth", "()I", func(thread *vm.Thread, args []interface{}) error {
		depth := int32(0)
		if traces := args[0].(*vm.Instance).AsThrowable(); traces != nil {
			depth = int32(len(traces))
		}

		thread.CurrentFrame().PushOperand(depth)
		return nil
	})

	vm.NativeMethods.Register(class, "getStackTraceElement", "(I)Ljava/lang/StackTraceElement;", func(thread *vm.Thread, args []interface{}) error {
		traces := args[0].(*vm.Instance).AsThrowable()
		index := args[1].(int32)

		if index < 0 || index >= int32(len(traces)) {
			return fmt.Errorf("index out of bounds")
		}

		thread.CurrentFrame().PushOperand(traces[index].ToJava(thread.VM()))
		return nil
	})
}
