package security

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func AccessControllerGetStackAccessControlContext(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(nil)
	return nil
}

func AccessControllerDoPrivilegedL(thread *vm.Thread, args []interface{}) error {
	action, ok := args[0].(*vm.Instance)
	if !ok {
		return fmt.Errorf("arg of AccessController.doPrivileged is NOT instance")
	}

	runClass, runMethod := action.Class().ResolveMethod("run", "()Ljava/lang/Object;")
	thread.CurrentFrame().PushOperand(action)
	return thread.ExecMethod(runClass, runMethod)
}
