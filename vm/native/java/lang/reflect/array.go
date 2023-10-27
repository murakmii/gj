package reflect

import "github.com/murakmii/gj/vm"

func ArrayNewArray(thread *vm.Thread, args []interface{}) error {
	compType := args[0].(*vm.Instance)
	size := args[1].(int32)
	array, _ := vm.NewArray(thread.VM(), "["+*(compType.VMData().(*string)), int(size))

	thread.CurrentFrame().PushOperand(array)
	return nil
}
