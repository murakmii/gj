package misc

import "github.com/murakmii/gj/vm"

func URLClassPathGetLookupCacheURLs(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(nil)
	return nil
}
