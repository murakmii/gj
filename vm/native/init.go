package native

import (
	"github.com/murakmii/gj/vm"
	"github.com/murakmii/gj/vm/native/java/lang"
)

func init() {
	vm.RegisterNativeMethod("java/lang/Class/desiredAssertionStatus0", lang.ClassDesiredAssertionStatus0)
	vm.RegisterNativeMethod("java/lang/Class/registerNatives", lang.ClassRegisterNatives)

	vm.RegisterNativeMethod("java/lang/Object/registerNatives", lang.ObjectRegisterNatives)
}
