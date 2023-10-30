package atomic

import "github.com/murakmii/gojiai/vm"

func init() {
	class := "java/util/concurrent/atomic/AtomicLong"

	vm.NativeMethods.Register(class, "VMSupportsCS8", "()Z", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(0))
		return nil
	})
}
