package security

import (
	"github.com/murakmii/gj/vm"
)

func AccessControllerGetStackAccessControlContext(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(nil)
	return nil
}
