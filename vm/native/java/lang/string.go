package lang

import (
	"github.com/murakmii/gj/vm"
)

func StringIntern(thread *vm.Thread, args []interface{}) error {
	gs := args[0].(*vm.Instance).AsString()

	interned, err := thread.VM().JavaString(thread, &gs)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(interned)
	return nil
}
