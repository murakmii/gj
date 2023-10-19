package reflect

import "github.com/murakmii/gj/vm"

func ReflectionGetCallerClass(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(nil)
	return nil
}
