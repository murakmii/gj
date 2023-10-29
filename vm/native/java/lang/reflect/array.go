package reflect

import "github.com/murakmii/gj/vm"

func init() {
	class := "java/lang/reflect/Array"

	vm.NativeMethods.Register(class, "newArray", "(Ljava/lang/Class;I)Ljava/lang/Object;", func(thread *vm.Thread, args []interface{}) error {
		compType := args[0].(*vm.Instance).AsClass().File().ThisClass()
		size := args[1].(int32)
		array, _ := vm.NewArray(thread.VM(), "["+compType, int(size))

		thread.CurrentFrame().PushOperand(array)
		return nil
	})
}
