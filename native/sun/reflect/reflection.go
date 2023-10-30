package reflect

import (
	"github.com/murakmii/gojiai/vm"
)

func init() {
	_class := "sun/reflect/Reflection"

	vm.NativeMethods.Register(_class, "getCallerClass", "()Ljava/lang/Class;", func(thread *vm.Thread, args []interface{}) error {
		callerClassName := thread.InvokerFrame().CurrentClass().File().ThisClass()
		class, err := thread.VM().Class(callerClassName, nil)
		if err != nil {
			return err
		}

		thread.CurrentFrame().PushOperand(class.Java())
		return nil
	})

	vm.NativeMethods.Register(_class, "getClassAccessFlags", "(Ljava/lang/Class;)I", func(thread *vm.Thread, args []interface{}) error {
		class := args[0].(*vm.Instance).AsClass()
		thread.CurrentFrame().PushOperand(int32(class.File().AccessFlag()))
		return nil
	})
}
