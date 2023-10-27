package atomic

import "github.com/murakmii/gj/vm"

func AtomicLongVMSupportsCS8(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(int32(0))
	return nil
}
