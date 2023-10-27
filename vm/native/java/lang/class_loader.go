package lang

import (
	"github.com/murakmii/gj/vm"
	"strings"
)

func ClassLoaderFindLoadedClass0(thread *vm.Thread, args []interface{}) error {
	className := strings.ReplaceAll(args[1].(*vm.Instance).GetCharArrayField("value"), ".", "/")
	class, err := thread.VM().Class(className, thread)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			thread.CurrentFrame().PushOperand(nil)
			return nil
		}
		return err
	}

	thread.CurrentFrame().PushOperand(class.Java())
	return nil
}
