package lang

import (
	"github.com/murakmii/gj/vm"
)

func init() {
	class := "java/lang/String"

	vm.NativeMethods.Register(class, "intern", "()Ljava/lang/String;", func(thread *vm.Thread, args []interface{}) error {
		gs := args[0].(*vm.Instance).AsString()

		interned, err := thread.VM().JavaString(thread, &gs)
		if err != nil {
			return err
		}

		thread.CurrentFrame().PushOperand(interned)
		return nil
	})
}
