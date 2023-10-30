package lang

import (
	"github.com/murakmii/gojiai/vm"
	"runtime"
)

func init() {
	class := "java/lang/Runtime"

	vm.NativeMethods.Register(class, "availableProcessors", "()I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(runtime.NumCPU()))
		return nil
	})
}
