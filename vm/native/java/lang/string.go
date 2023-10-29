package lang

import (
	"github.com/murakmii/gojiai/vm"
)

func init() {
	class := "java/lang/String"

	vm.NativeMethods.Register(class, "intern", "()Ljava/lang/String;", func(thread *vm.Thread, args []interface{}) error {
		gs := args[0].(*vm.Instance).AsString()
		thread.CurrentFrame().PushOperand(thread.VM().JavaString(gs))
		return nil
	})
}
